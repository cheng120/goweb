package id

import "github.com/rs/xid"

type IDService struct {
}

func NewIDService(params ...interface{}) (interface{}, error) {
	return &IDService{}, nil
}

func (s *IDService) NewID() string {
	return xid.New().String()
}