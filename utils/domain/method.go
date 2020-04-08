package domain

import "time"

type Method struct {
	Method     string
	Url        string
	Headers    map[string]string
	ParamsType string
	Body       []byte
	Params     map[string]string
	Timeout    time.Duration
}
