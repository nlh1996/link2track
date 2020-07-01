package model

import "sync"

// SpanMap .
var SpanMap map[string]Spans

// ErrTid .
var ErrTid map[string]string

// Stream .
var Stream chan Span

// Mux 互斥锁
var Mux sync.Mutex

// Span .
type Span struct {
	Tid  string
	Time string
	Data string
}

// Result .
var Result map[string]string

// EndSign 下载结束信号
var EndSign int

// Init .
func Init() {
	Result = make(map[string]string, 10000)
	ErrTid = make(map[string]string)
	Stream = make(chan Span, 31000)
	SpanMap = make(map[string]Spans, 10000)
}

// Spans .
type Spans []Span

// Len
func (s Spans) Len() int {
	return len(s)
}

// Less 排序
func (s Spans) Less(i, j int) bool {
	return s[i].Time < s[j].Time
}

// Swap 
func (s Spans) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
