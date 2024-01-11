package server

import "errors"

var (
	// ValueNotFoundErr is used when a requested value not found.
	ValueNotFoundErr = errors.New("value not found")
	// UnsupportedCommandErr is used when the server receives message with invalid command name.
	UnsupportedCommandErr = errors.New("command is not supported")
)
