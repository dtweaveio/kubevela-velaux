package api

import (
	"fmt"
	"k8s.io/apimachinery/pkg/watch"
	"net/http"
)

// SSEWriter 是一个封装了 http.ResponseWriter 的结构体
type SSEWriter struct {
	writer http.ResponseWriter
}

// NewSSEWriter 创建一个新的 SSEWriter
func NewSSEWriter(w http.ResponseWriter) *SSEWriter {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Transfer-Encoding", "chunked")
	return &SSEWriter{writer: w}
}

// SendMessage 发送 SSE 消息
func (sse *SSEWriter) SendMessage(data watch.EventType) error {
	message := fmt.Sprintf("data: %s\n\n", data)
	_, err := sse.writer.Write([]byte(message))
	if err != nil {
		return err
	}
	sse.writer.(http.Flusher).Flush()
	return nil
}
