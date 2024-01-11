package server

import (
	"bytes"
	"fmt"
)

type Item struct {
	key   string
	value string
}

func (i Item) Serialize() []byte {
	var buf bytes.Buffer

	_, _ = fmt.Fprintf(&buf, "%q:%q\n", i.key, i.value)
	return buf.Bytes()
}
