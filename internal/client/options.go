package client

import (
	"time"

	"github.com/zy3101176/gpt-4o-image-go/pkg/models"
)

// Options 客户端配置选项
type Options struct {
	APIKey     string
	APIURL     string
	RateLimit  time.Duration
	MaxWorkers int
}

// defaultOptions 返回默认配置
func defaultOptions() *Options {
	return &Options{
		APIURL:     "https://api.tu-zi.com/v1/chat/completions",
		RateLimit:  500 * time.Millisecond,
		MaxWorkers: 5,
	}
}

// WithAPIKey 设置API密钥
func WithAPIKey(apiKey string) models.Option[Options] {
	return func(o *Options) {
		o.APIKey = apiKey
	}
}

// WithAPIURL 设置API URL
func WithAPIURL(url string) models.Option[Options] {
	return func(o *Options) {
		o.APIURL = url
	}
}

// WithRateLimit 设置请求速率限制
func WithRateLimit(rateLimit time.Duration) models.Option[Options] {
	return func(o *Options) {
		o.RateLimit = rateLimit
	}
}

// WithMaxWorkers 设置最大工作协程数
func WithMaxWorkers(maxWorkers int) models.Option[Options] {
	return func(o *Options) {
		o.MaxWorkers = maxWorkers
	}
}
