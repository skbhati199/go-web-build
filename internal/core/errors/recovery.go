package errors

import (
	"fmt"
	"runtime/debug"
	"sync"
)

type RecoveryManager struct {
	handlers map[string]RecoveryHandler
	mu       sync.RWMutex
}

type RecoveryHandler func(err interface{}) error

func NewRecoveryManager() *RecoveryManager {
	return &RecoveryManager{
		handlers: make(map[string]RecoveryHandler),
	}
}

func (rm *RecoveryManager) RegisterHandler(name string, handler RecoveryHandler) {
	rm.mu.Lock()
	defer rm.mu.Unlock()
	rm.handlers[name] = handler
}

func (rm *RecoveryManager) Recover(name string) {
	if r := recover(); r != nil {
		rm.mu.RLock()
		handler, exists := rm.handlers[name]
		rm.mu.RUnlock()

		if exists {
			if err := handler(r); err != nil {
				debug.PrintStack()
				panic(fmt.Sprintf("recovery handler failed: %v", err))
			}
		} else {
			debug.PrintStack()
			panic(fmt.Sprintf("no recovery handler for: %s", name))
		}
	}
}
