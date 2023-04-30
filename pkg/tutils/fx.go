package tutils

import "go.uber.org/fx"

type EmtpyLifecycle struct {
	Hooks []fx.Hook
}

func (e *EmtpyLifecycle) Append(h fx.Hook) {
	e.Hooks = append(e.Hooks, h)
}

func NewEmtpyLifecycle() *EmtpyLifecycle {
	return &EmtpyLifecycle{}
}
