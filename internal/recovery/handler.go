package recovery

import (
	"fmt"
	"runtime/debug"

	"github.com/sonukumar/go-web-build/internal/logger"
)

type RecoveryHandler struct {
	debug bool
}

func NewRecoveryHandler(debug bool) *RecoveryHandler {
	return &RecoveryHandler{
		debug: debug,
	}
}

func (h *RecoveryHandler) Recover(err interface{}) error {
	stack := debug.Stack()

	if h.debug {
		logger.Error(fmt.Errorf("%v", err), "Panic recovered\nStack trace:")
		logger.Debug(string(stack))
	} else {
		logger.Error(fmt.Errorf("%v", err), "An unexpected error occurred")
	}

	return fmt.Errorf("recovered from panic: %v", err)
}

func (h *RecoveryHandler) WrapHandler(fn func() error) func() error {
	return func() error {
		defer func() {
			if r := recover(); r != nil {
				h.Recover(r)
			}
		}()
		return fn()
	}
}
