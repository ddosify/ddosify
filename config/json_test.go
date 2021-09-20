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

package config

import (
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"ddosify.com/hammer/core/types"
)

func TestCreateHammerDefaultValues(t *testing.T) {
	t.Parallel()
	jsonReader, _ := NewConfigReader("config_testdata/config_empty.json", ConfigTypeJson)
	expectedHammer := types.Hammer{
		TotalReqCount:     types.DefaultReqCount,
		LoadType:          types.DefaultLoadType,
		TestDuration:      types.DefaultDuration,
		ReportDestination: types.DefaultOutputType,
		Scenario: types.Scenario{
			Scenario: []types.ScenarioItem{{
				ID:       1,
				URL:      strings.ToLower(types.DefaultProtocol) + "://test.com",
				Protocol: types.DefaultProtocol,
				Method:   types.DefaultMethod,
				Timeout:  types.DefaultTimeout,
			}},
		},
		Proxy: types.Proxy{
			Strategy: "single",
		},
	}

	h, err := jsonReader.CreateHammer()

	if err != nil {
		t.Errorf("TestCreateHammerDefaultValues error occured: %v", err)
	}

	if !reflect.DeepEqual(expectedHammer, h) {
		t.Errorf("Expected: %v, Found: %v", expectedHammer, h)
	}
}

func TestCreateHammer(t *testing.T) {
	t.Parallel()
	jsonReader, _ := NewConfigReader("config_testdata/config.json", ConfigTypeJson)
	addr, _ := url.Parse("http://lothygyo:vezpdbtjk731@209.127.191.180:9279")
	expectedHammer := types.Hammer{
		TotalReqCount:     1555,
		LoadType:          types.LoadTypeWaved,
		TestDuration:      21,
		ReportDestination: types.OutputTypeTimescale,
		Scenario: types.Scenario{
			Scenario: []types.ScenarioItem{
				{
					ID:       1,
					URL:      "https://app.servdown.com/accounts/login/?next=/",
					Protocol: types.ProtocolHTTPS,
					Method:   http.MethodGet,
					Timeout:  3,
					Payload:  "payload str",
					Custom: map[string]interface{}{
						"keep-alive": true,
					},
				},
				{
					ID:       2,
					URL:      "http://test.com",
					Protocol: types.ProtocolHTTP,
					Method:   http.MethodPut,
					Timeout:  2,
					Headers: map[string]string{
						"ContenType":    "application/xml",
						"X-ddosify-key": "ajkndalnasd",
					},
				},
			},
		},
		Proxy: types.Proxy{
			Strategy: "single",
			Addr:     addr,
		},
	}

	h, err := jsonReader.CreateHammer()

	if err != nil {
		t.Errorf("TestCreateHammerPro error occured: %v", err)
	}

	if !reflect.DeepEqual(expectedHammer, h) {
		t.Errorf("Expected: %v, Found: %v", expectedHammer, h)
	}
}

func TestCreateHammerPayload(t *testing.T) {
	t.Parallel()
	jsonReader, _ := NewConfigReader("config_testdata/config_payload.json", ConfigTypeJson)
	expectedPayloads := []string{"payload from string", "Payloaf from file."}
	h, err := jsonReader.CreateHammer()

	if err != nil {
		t.Errorf("TestCreateHammerPayload error occured: %v", err)
	}

	steps := h.Scenario.Scenario

	if steps[0].Payload != expectedPayloads[0] {
		t.Errorf("Expected: %v, Found: %v", expectedPayloads[0], steps[0].Payload)
	}

	if steps[1].Payload != expectedPayloads[1] {
		t.Errorf("Expected: %v, Found: %v", expectedPayloads[1], steps[1].Payload)
	}
}

func TestCreateHammerAuth(t *testing.T) {
	t.Parallel()
	jsonReader, _ := NewConfigReader("config_testdata/config_auth.json", ConfigTypeJson)
	expectedAuths := []types.Auth{
		{
			Type:     types.AuthHttpBasic,
			Username: "kursat",
			Password: "12345",
		},
		{}}

	h, err := jsonReader.CreateHammer()
	if err != nil {
		t.Errorf("TestCreateHammerAuth error occured: %v", err)
	}

	steps := h.Scenario.Scenario
	if steps[0].Auth != expectedAuths[0] {
		t.Errorf("Expected: %v, Found: %v", expectedAuths[0], steps[0].Auth)
	}

	if steps[1].Auth != expectedAuths[1] {
		t.Errorf("Expected: %v, Found: %v", expectedAuths[1], steps[1].Auth)
	}
}

func TestCreateHammerProtocol(t *testing.T) {
	t.Parallel()
	jsonReader, _ := NewConfigReader("config_testdata/config_protocol.json", ConfigTypeJson)
	expectedProtocols := []string{"HTTPS", "HTTP", types.DefaultProtocol, "HTTP"}

	h, err := jsonReader.CreateHammer()
	if err != nil {
		t.Errorf("TestCreateHammerProtocol error occured: %v", err)
	}

	steps := h.Scenario.Scenario
	for i := 0; i < len(steps); i++ {
		if steps[i].Protocol != expectedProtocols[i] {
			t.Errorf("1: Expected: %v, Found: %v", expectedProtocols[i], steps[0].Protocol)
		}

		url, err := url.Parse(steps[i].URL)
		if err != nil {
			t.Errorf("TestCreateHammerProtocol-SchemeCheck error occured: %v", err)
		}

		if strings.ToUpper(url.Scheme) != expectedProtocols[i] {
			t.Errorf("2: Expected: %v, Found: %v", expectedProtocols[i], url.Scheme)
		}
	}
}

func TestCreateHammerInvalidTarget(t *testing.T) {
	t.Parallel()
	jsonReader, _ := NewConfigReader("config_testdata/config_invalid_target.json", ConfigTypeJson)

	_, err := jsonReader.CreateHammer()
	if err == nil {
		t.Errorf("TestCreateHammerProtocol error occured")
	}
}