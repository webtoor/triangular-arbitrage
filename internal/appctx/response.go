package appctx

import (
	"encoding/json"
	"sync"
)

var (
	rsp    *Response
	oneRsp sync.Once
)

// Response presentation contract object
type Response struct {
	Code        int    `json:"code,omitempty"`
	Status      string `json:"status,omitempty"`
	Message     any    `json:"message,omitempty"`
	Errors      any    `json:"errors,omitempty"`
	Data        any    `json:"data,omitempty"`
	rawResponse any    `json:"-"`
}

// Generate setter message
func (r *Response) Generate() *Response {
	return r
}

// WithCode setter response var name
func (r *Response) WithCode(c int) *Response {
	r.Code = c
	return r
}

// WithData setter data response
func (r *Response) WithData(v any) *Response {
	r.Data = v
	return r
}

// WithError setter error messages
func (r *Response) WithError(v any) *Response {
	r.Errors = v
	return r
}

// WithMessage setter custom message response
func (r *Response) WithMessage(v any) *Response {
	if v != nil {
		r.Message = v
	}

	return r
}

// WithRawResponse setter raw response
func (r *Response) WithRawResponse(v any) *Response {
	r.rawResponse = v
	return r
}

// Byte cast response to byte
func (r *Response) Byte() []byte {
	if r.Code == 0 || r.Message == nil {
		r.Generate()
	}

	b, _ := json.Marshal(r)
	return b
}

// NewResponse initialize response
func NewResponse() *Response {
	oneRsp.Do(func() {
		rsp = &Response{}
	})

	// clone response
	x := *rsp

	return &x
}

func (r *Response) RawResponse() any {
	return r.rawResponse
}
