package loadgen

import (
	"Go-Demo/helper/log"
	"Go-Demo/loadgen/lib"
	"bytes"
	"context"
	"fmt"
	"math"
	"sync/atomic"
	"time"
)

// 日志记录器。
var logger = log.DLogger()

// myGenerator 代表载荷发生器的实现类型。
type myGenerator struct {
	caller lib.Caller // 调用器
	timeoutNS time.Duration // 处理超市时间,单位:纳秒
	lps uint32 // 每秒载荷量
	durationNS time.Duration // 负载持续时间,单位:纳秒
	concurrency uint32 // 载荷并发量
	tickets lib.GoTickets // Goroutine票池
	ctx context.Context // 上下文
	cancelFunc context.CancelFunc // 取消函数
	callCount int64 // 调用计数
	status uint32 // 状态
	resultCh chan *lib.CallResult // 调用结果通道
}

// NewGenerator 会创建一个载荷发生器.
func NewGenerator(pset ParamSet) (lib.Generator, error) {
	logger.Infoln("New a load generator...")
	if err := pset.Check(); err != nil {
		return nil, err
	}
	gen := &myGenerator{
		caller:     pset.Caller,
		timeoutNS:  pset.TimeoutNS,
		lps:        pset.LPS,
		durationNS: pset.DurationNS,
		status:     lib.STATUS_ORIGINAL,
		resultCh:   pset.ResultCh,
	}
	if err := gen.init(); err != nil {
		return nil, err
	}
	return gen, nil
}

// 初始化载荷发生器
func (gen *myGenerator) init() error {
	var buf bytes.Buffer
	buf.WriteString("Initializing the load generator...")
	// 载荷的并发量 ≈ 载荷的响应超时时间 / 载荷的发送间隔时间
	var total64 = int64(gen.timeoutNS)/int64(1e9/gen.lps) + 1
	if total64 > math.MaxInt32 {
		total64 = math.MaxInt32
	}
	gen.concurrency = uint32(total64)
	tickets, err := lib.NewGoTickets(gen.concurrency)
	if err != nil {
		return err
	}
	gen.tickets = tickets

	buf.WriteString(fmt.Sprintf("Done. (concurrency=%d)", gen.concurrency))
	logger.Infoln(buf.String())
	return nil
}

// Start 会启动载荷发生器
func (gen *myGenerator) Start() bool {
	logger.Infoln("Starting load generator...")
	// 检查是否具备可启动的状态, 顺便设置状态为正在启动
	if !atomic.CompareAndSwapUint32(
		&gen.status, lib.STATUS_ORIGINAL, lib.STATUS_STARTING) {
		if !atomic.CompareAndSwapUint32(&gen.status, lib.STATUS_STOPPED, lib.STATUS_STARTING) {
			return false
		}
	}

	// 设定节流阀.
	var throttle <- chan time.Time
	if gen.lps > 0 {
		interval := time.Duration(1e9 / gen.lps)
		logger.Infof("Setting throttle (%v)...", interval)
		throttle = time.Tick(interval)
	}

	// 初始化上下文和取消函数.
	gen.ctx, gen.cancelFunc = context.WithTimeout(
		context.Background(), gen.durationNS)

	// 初始化调用计数.
	gen.callCount = 0

	// 设置状态为已启动.
	atomic.StoreUint32(&gen.status, lib.STATUS_STARTED)

	go func() {
		// 生成并发送载荷
		logger.Infoln("Generating loads...")
		gen.genLoad(throttle)
		logger.Info("Stopped. (call count: %d", gen.callCount)
	}()
	return true
}

func (gen *myGenerator) Status() uint32 {
	return atomic.LoadUint32(&gen.status)
}

func (gen *myGenerator) CallCount() int64 {
	return atomic.LoadInt64(&gen.callCount)
}