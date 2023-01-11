package x

import "net/http"

type Response interface {
	IsOk() bool
	SetTraceId(traceId string)
	GetData() H
	GetMessage() string
	SetTime(time string)
}

func NewResponse(code int, msg string, data H) Response {
	return &DefaultResponse{Code: code, Msg: msg, Data: data}
}

type DefaultResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data H      `json:"data"`

	TraceID string `json:"trace_id"`
	Time    string `json:"time"`
}

func (h *DefaultResponse) SetTime(time string) {
	h.Time = time
}

func (h *DefaultResponse) GetData() H {
	return h.Data
}

func (h *DefaultResponse) GetMessage() string {
	return h.Msg
}

func (h *DefaultResponse) SetTraceId(traceId string) {
	h.TraceID = traceId
}

func (h *DefaultResponse) IsOk() bool {
	return h.Code == http.StatusOK
}

