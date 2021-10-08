package scheduler

import "Go-Demo/webcrawler/module"

// SchedSummary 代表调度器摘要的接口类型。
type SchedSummary interface {
	// Struct 用于获得摘要信息的结构化形式。
	Struct() module.SummaryStruct
	// String 用于获得摘要信息的字符串形式。
	String() string
}

// newSchedSummary 用于创建一个调度器摘要实例。
func newSchedSummary(
	requestArgs RequestArgs,
	dataArgs DataArgs,
	moduleArgs ModuleArgs,
	sched *myScheduler) SchedSummary {
	if sched == nil {
		return nil
	}
	return &mySchedSummary{
		requestArgs: requestArgs,
		dataArgs: dataArgs,
		moduleArgs: moduleArgs,
		sched: sched,
	}
}

// mySchedSummary 代表调度器摘要的实现类型。
type mySchedSummary struct {
	// requestArgs 代表请求相关的参数。
	requestArgs RequestArgs
	// dataArgs 代表数据相关参数的容器实例。
	dataArgs DataArgs
	// moduleArgs 代表组件相关参数的容器实例。
	moduleArgs ModuleArgs
	// maxDepth 代表爬取的最大深度。
	maxDepth uint32
	// sched 代表调度器实例。
	sched *myScheduler
}

// SummaryStruct 代表调度器摘要的结构。
type SummaryStruct struct {
	RequestArgs     RequestArgs             `json:"request_args"`
	DataArgs        DataArgs                `json:"data_args"`
	ModuleArgs      ModuleArgsSummary       `json:"module_args"`
	Status          string                  `json:"status"`
	Downloaders     []module.SummaryStruct  `json:"downloaders"`
	Analyzers       []module.SummaryStruct  `json:"analyzers"`
	Pipelines       []module.SummaryStruct  `json:"pipelines"`
	ReqBufferPool   BufferPoolSummaryStruct `json:"request_buffer_pool"`
	RespBufferPool  BufferPoolSummaryStruct `json:"response_buffer_pool"`
	ItemBufferPool  BufferPoolSummaryStruct `json:"item_buffer_pool"`
	ErrorBufferPool BufferPoolSummaryStruct `json:"error_buffer_pool"`
	NumURL          uint64                  `json:"url_number"`
}