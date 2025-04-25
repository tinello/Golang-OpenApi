package sys_actions

import (
	sys_domain "github.com/tinello/golang-openapi/core/system/domain"
)

//go:generate mockgen -package mocks -destination ../../mocks/MockGetServiceInfo.go . GetServiceInfo

type GetServiceInfo interface {
	Execute() *sys_domain.ServiceInfo
}

func NewGetServiceInfo() *getServiceInfo {
	return &getServiceInfo{}
}

type getServiceInfo struct {
}

func (s *getServiceInfo) Execute() *sys_domain.ServiceInfo {
	return &sys_domain.ServiceInfo{
		Error: nil,
	}
}
