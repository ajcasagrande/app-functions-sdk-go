//
// Copyright (c) 2019 Intel Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

package webserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/edgexfoundry/app-functions-sdk-go/internal"
	"github.com/edgexfoundry/app-functions-sdk-go/internal/common"
	"github.com/edgexfoundry/app-functions-sdk-go/internal/telemetry"
	"github.com/edgexfoundry/go-mod-core-contracts/clients"
	"github.com/edgexfoundry/go-mod-core-contracts/clients/logger"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var logClient logger.LoggingClient
var config *common.ConfigurationStruct

func init() {
	logClient = logger.NewClient("app_functions_sdk_go", false, "./test.log", "DEBUG")
	config = &common.ConfigurationStruct{}
}

func TestConfigureAndPingRoute(t *testing.T) {

	webserver := NewWebServer(config, logClient, mux.NewRouter())
	webserver.ConfigureStandardRoutes()

	req, _ := http.NewRequest("GET", clients.ApiPingRoute, nil)
	rr := httptest.NewRecorder()
	webserver.router.ServeHTTP(rr, req)

	body := rr.Body.String()
	assert.Equal(t, "pong", body)

}

func TestConfigureAndVersionRoute(t *testing.T) {

	webserver := NewWebServer(config, logClient, mux.NewRouter())
	webserver.ConfigureStandardRoutes()

	req, _ := http.NewRequest("GET", clients.ApiVersionRoute, nil)
	rr := httptest.NewRecorder()
	webserver.router.ServeHTTP(rr, req)

	body := rr.Body.String()
	assert.Equal(t, "{\"version\":\"0.0.0\",\"sdk_version\":\"0.0.0\"}\n", body)

}
func TestConfigureAndConfigRoute(t *testing.T) {

	webserver := NewWebServer(config, logClient, mux.NewRouter())
	webserver.ConfigureStandardRoutes()

	req, _ := http.NewRequest("GET", clients.ApiConfigRoute, nil)
	rr := httptest.NewRecorder()
	webserver.router.ServeHTTP(rr, req)

	expected := `{"Writable":{"LogLevel":"","Pipeline":{"ExecutionOrder":"","UseTargetTypeOfByteArray":false,"Functions":null},"StoreAndForward":{"Enabled":false,"RetryInterval":"","MaxRetryCount":0}},"Logging":{"EnableRemote":false,"File":""},"Registry":{"Host":"","Port":0,"Type":""},"Service":{"BootTimeout":"","CheckInterval":"","ClientMonitor":"","Host":"","Port":0,"Protocol":"","StartupMsg":"","ReadMaxLimit":0,"Timeout":""},"MessageBus":{"PublishHost":{"Host":"","Port":0,"Protocol":""},"SubscribeHost":{"Host":"","Port":0,"Protocol":""},"Type":"","Optional":null},"Binding":{"Type":"","SubscribeTopic":"","PublishTopic":""},"ApplicationSettings":null,"Clients":null,"Database":{"Type":"","Host":"","Port":0,"Timeout":"","Username":"","Password":"","MaxIdle":0,"BatchSize":0},"SecretStore":{"Host":"","Port":0,"Path":"","Protocol":"","Namespace":"","RootCaCertPath":"","ServerName":"","Authentication":{"AuthType":"","AuthToken":""},"TokenFile":""}}` + "\n"

	body := rr.Body.String()
	assert.Equal(t, expected, body)
}

func TestConfigureAndMetricsRoute(t *testing.T) {
	webserver := NewWebServer(config, logClient, mux.NewRouter())
	webserver.ConfigureStandardRoutes()

	req, _ := http.NewRequest("GET", clients.ApiMetricsRoute, nil)
	rr := httptest.NewRecorder()
	webserver.router.ServeHTTP(rr, req)

	body := rr.Body.String()
	metrics := telemetry.SystemUsage{}
	json.Unmarshal([]byte(body), &metrics)
	assert.NotNil(t, body, "Metrics not populated")
	assert.NotZero(t, metrics.Memory.Alloc, "Expected Alloc value of metrics to be non-zero")
	assert.NotZero(t, metrics.Memory.Frees, "Expected Frees value of metrics to be non-zero")
	assert.NotZero(t, metrics.Memory.LiveObjects, "Expected LiveObjects value of metrics to be non-zero")
	assert.NotZero(t, metrics.Memory.Mallocs, "Expected Mallocs value of metrics to be non-zero")
	assert.NotZero(t, metrics.Memory.Sys, "Expected Sys value of metrics to be non-zero")
	assert.NotZero(t, metrics.Memory.TotalAlloc, "Expected TotalAlloc value of metrics to be non-zero")
	assert.NotNil(t, metrics.CpuBusyAvg, "Expected CpuBusyAvg value of metrics to be not nil")
}

func TestSetupTriggerRoute(t *testing.T) {
	webserver := NewWebServer(config, logClient, mux.NewRouter())

	handlerFunctionNotCalled := true
	handler := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
		handlerFunctionNotCalled = false
	}

	webserver.SetupTriggerRoute(handler)

	req, _ := http.NewRequest("GET", internal.ApiTriggerRoute, nil)
	rr := httptest.NewRecorder()
	webserver.router.ServeHTTP(rr, req)

	body := rr.Body.String()

	assert.Equal(t, "test", body)
	assert.False(t, handlerFunctionNotCalled, "expected handler function to be called")

}
