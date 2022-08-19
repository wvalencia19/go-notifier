package config

import "time"

type Notifier struct {
	Interval         time.Duration `default:"50ms"`
	MaxAllowedErrors int           `default:"4"`
	HTTPClient       HTTPClient
	LogLevel         string `default:"info"`
	BufferSize       int    `default:"0"` // Default size for buffer of the messages channel
}
type HTTPClient struct {
	HTTPRequestTimeout  time.Duration `default:"1s"`
	MaxIdleConns        int           `default:"100"`
	MaxIdleConnsPerHost int           `default:"100"`
}
