package sys_domain

type ServiceInfo struct {
	Error error
}

func (r ServiceInfo) Healthy() bool {
	return r.Error == nil
}
