package proxy

import (
	"net/url"
	"strings"

	"ddosify.com/hammer/core/types"
)

type ProxyService interface {
	init(types.Proxy, int)
	GetNewProxy() *url.URL
	ReportProxy(addr *url.URL, reason string)
}

func NewProxyService(p types.Proxy, reqCount int) (service ProxyService, err error) {
	if strings.EqualFold(p.Strategy, "single") {
		service = &singleProxyStrategy{}
	}
	service.init(p, reqCount)

	return service, nil
}
