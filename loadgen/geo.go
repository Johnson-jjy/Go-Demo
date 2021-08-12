package loadgen

import (
	"Go-Demo/loadgen/lib"
	"Go-Demo/loadgen/lib/log"
	"context"
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
	tickets lib.Gotickets // Goroutine票池
	ctx context.Context // 上下文
	cancelFunc context.CancelFunc // 取消函数
	callCount int64 // 调用计数
	status uint32 // 状态
	result chan *lib.CallResult // 调用结果通道
}

// NewGenerator 会创建一个载荷发生器.
func NewGenerator(pset ParamSet) (lib.Generator, error) {
	logger.Infoln("New a load generator...")
	if err := pset.Check(); err != nil {
		return nil, err
	}
}