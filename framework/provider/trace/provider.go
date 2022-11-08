package trace

import (
	"goweb/framework"
	"goweb/framework/contract"
)

type TraceProvider struct {
	c framework.Container
}

// Register registe a new function for make a service instance
func (provider *TraceProvider) Register(c framework.Container) framework.NewInstance {
	return NewTraceService
}

// Boot will called when the service instantiate
func (provider *TraceProvider) Boot(c framework.Container) error {
	provider.c = c
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *TraceProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *TraceProvider) Params(c framework.Container) []interface{} {
	return []interface{}{provider.c}
}

/// Name define the name for this service
func (provider *TraceProvider) Name() string {
	return contract.TraceKey
}