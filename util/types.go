package util

import "context"

type ContextWrapper interface {
	GetContext() context.Context
}
