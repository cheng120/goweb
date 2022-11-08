package id

import (
	"goweb/framework"
	"goweb/framework/contract"
)

type IDProvider struct {
}

// Register registe a new function for make a service instance
func (provider *IDProvider) Register(c framework.Container) framework.NewInstance {
	return NewIDService
}

// Boot will called when the service instantiate
func (provider *IDProvider) Boot(c framework.Container) error {
	return nil
}

// IsDefer define whether the service instantiate when first make or register
func (provider *IDProvider) IsDefer() bool {
	return false
}

// Params define the necessary params for NewInstance
func (provider *IDProvider) Params(c framework.Container) []interface{} {
	return []interface{}{}
}

/// Name define the name for this service
func (provider *IDProvider) Name() string {
	return contract.IDKey
}