package httpclient

import (
	"time"

	"github.com/avast/retry-go/v4"
)

type TRetryDelayType int

const (
	RetryDelayTypeConstant    TRetryDelayType = iota // 固定延遲時間
	RetryDelayTypeExponential                        // 指數退避延遲
)

type RetryConfig struct {
	UntilToSuccess bool            // 直到成功為止(設定這個，就不管Attempts)
	Attempts       uint            // 嘗試次數
	Delay          time.Duration   // 等待時間
	MaxDelay       time.Duration   // 最大等待時間
	DelayType      TRetryDelayType // 延遲類型：固定延遲或指數退避
	LastErrorOnly  bool            // 只顯示最後一個錯誤
}

func (receiver *RetryConfig) ToRetryOpts() []retry.Option {
	retryOpts := make([]retry.Option, 0)
	retryOpts = append(retryOpts, retry.Attempts(receiver.Attempts))
	if receiver.UntilToSuccess {
		retryOpts = append(retryOpts, retry.UntilSucceeded())
	}
	if receiver.Delay > 0 {
		retryOpts = append(retryOpts, retry.Delay(receiver.Delay))
	}
	if receiver.MaxDelay > 0 {
		retryOpts = append(retryOpts, retry.MaxDelay(receiver.MaxDelay))
	}
	switch receiver.DelayType {
	case RetryDelayTypeConstant:
		retryOpts = append(retryOpts, retry.DelayType(retry.FixedDelay))
	case RetryDelayTypeExponential:
		retryOpts = append(retryOpts, retry.DelayType(retry.BackOffDelay))
	}
	retryOpts = append(retryOpts, retry.LastErrorOnly(receiver.LastErrorOnly))

	return retryOpts
}
