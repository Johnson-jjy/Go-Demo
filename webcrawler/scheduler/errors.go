package scheduler

import "Go-Demo/webcrawler/errors"

// genError 用于生成爬虫错误值。
func genError(errMsg string) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_SCHEDULER,
		errMsg)
}

// genErrorByError 用于基于给定的错误值生成爬虫错误值。
func genErrorByError(err error) error {
	return errors.NewCrawlerError(errors.ERROR_TYPE_SCHEDULER,
		err.Error())
}

// genParameterError 用于生成爬虫参数错误值。
func genParameterError(errMsg string) error {
	return errors.NewCrawlerErrorBy(errors.ERROR_TYPE_SCHEDULER,
		errors.NewIllegalParameterError(errMsg))
}
