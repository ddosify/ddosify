/*
*
*	Ddosify - Load testing tool for any web system.
*   Copyright (C) 2021  Ddosify (https://ddosify.com)
*
*   This program is free software: you can redistribute it and/or modify
*   it under the terms of the GNU Affero General Public License as published
*   by the Free Software Foundation, either version 3 of the License, or
*   (at your option) any later version.
*
*   This program is distributed in the hope that it will be useful,
*   but WITHOUT ANY WARRANTY; without even the implied warranty of
*   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
*   GNU Affero General Public License for more details.
*
*   You should have received a copy of the GNU Affero General Public License
*   along with this program.  If not, see <https://www.gnu.org/licenses/>.
*
 */

package proxy

import (
	"fmt"
	"net/url"
	"strings"

	"ddosify.com/hammer/core/types"
)

type ProxyService interface {
	Init(types.Proxy) error
	GetAll() []*url.URL
	GetProxy() *url.URL
	ReportProxy(addr *url.URL, reason string) *url.URL
	GetProxyCountry(*url.URL) string
}

func NewProxyService(s string) (service ProxyService, err error) {
	if strings.EqualFold(s, "single") {
		service = &singleProxyStrategy{}
	} else {
		err = fmt.Errorf("unsupported proxy strategy")
	}

	return
}
