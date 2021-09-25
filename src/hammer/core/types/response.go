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

package types

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

//Equivalent to Scenario. Each Scenario has a Response after request is done.
type Response struct {
	// First request start time for the Scenario
	StartTime time.Time

	ProxyAddr     *url.URL
	ResponseItems []*ResponseItem
}

//Equivalent to ScenarioItem.
type ResponseItem struct {
	// ID of the ScenarioItem
	ScenarioItemID int16

	// Each request has a unique ID.
	RequestID uuid.UUID

	// Returned status code. Has different meaning for different protocols.
	StatusCode int

	// Time of the request call.
	RequestTime time.Time

	// Total duration. From request sending to full response recieving.
	Duration time.Duration

	// Response content length
	ContentLenth int64

	// Error occured at request time.
	Err RequestError

	// Protocol spesific metrics. For ex: DNSLookupDuration: 1s for HTTP
	Custom map[string]interface{}
}
