package tests

import (
	"bytes"
	"sync"
)

type MemLog struct {
	buf *bytes.Buffer
	sync.Mutex
}

func NewMemLog() *MemLog {
	return &MemLog{
		buf:   bytes.NewBuffer(nil),
		Mutex: sync.Mutex{},
	}
}

func (m *MemLog) String() string {
	m.Lock()
	defer m.Unlock()

	return m.buf.String()
}

func (m *MemLog) Write(p []byte) (n int, err error) {
	m.Lock()
	defer m.Unlock()

	return m.buf.Write(p)
}
