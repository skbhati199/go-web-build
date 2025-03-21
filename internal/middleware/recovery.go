package middleware

import "github.com/skbhati199/go-web-build/internal/core/errors"

type RecoveryMiddleware struct {
	manager *errors.RecoveryManager
}

func NewRecoveryMiddleware(manager *errors.RecoveryManager) *RecoveryMiddleware {
	return &RecoveryMiddleware{
		manager: manager,
	}
}

func (rm *RecoveryMiddleware) Wrap(name string, fn func() error) func() error {
	return func() error {
		defer rm.manager.Recover(name)
		return fn()
	}
}
