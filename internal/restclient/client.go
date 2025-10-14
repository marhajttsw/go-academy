package restclient

import (
	"time"

	"github.com/go-resty/resty/v2"
)

type Config struct {
	BaseURL string
	Timeout time.Duration
	Headers map[string]string
}

func New(cfg Config) *resty.Client {
	c := resty.New().
		SetBaseURL(cfg.BaseURL).
		SetTimeout(cfg.Timeout)

	for k, v := range cfg.Headers {
		c.SetHeader(k, v)
	}
	return c
}
