// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package components

import (
	"fmt"
	"net"
	"time"

	"github.com/GoogleCloudPlatform/api-proxy/tests/env/testdata"
	"github.com/golang/glog"
)

// Dynamic port allocation scheme
// To avoid port conflicts when setting up test env. Each test should use unique ports
// Each test has a unique test_name, its ports will be allocated based on that name

// All integration tests should be listed here to get their test ids
const (
	TestAsymmetricKeys uint16 = iota
	TestAsymmetricKeysTraces
	TestAuthJwksCache
	TestBackendAuthWithImdsIdToken
	TestBackendAuthWithIamIdToken
	TestDifferentOriginPreflightCors
	TestDifferentOriginSimpleCors
	TestDynamicRouting
	TestGRPC
	TestGrpcBackendPreflightCors
	TestGrpcBackendSimpleCors
	TestGRPCErrors
	TestGRPCFallback
	TestGRPCInteropMiniStress
	TestGRPCInterops
	TestGRPCJwt
	TestGRPCMetadata
	TestGRPCMinistress
	TestGRPCStreaming
	TestGRPCLongStreaming
	TestGRPCWeb
	TestHttp1Basic
	TestHttp1JWT
	TestManagedServiceConfig
	TestPreflightCorsWithBasicPreset
	TestPreflightRequestWithAllowCors
	TestReportGCPAttributes
	TestServiceControlALlHTTPMethod
	TestServiceControlAPIKeyCustomLocation
	TestServiceControlAPIKeyDefaultLocation
	TestServiceControlAPIKeyIpRestriction
	TestServiceControlAPIKeyRestriction
	TestServiceControlBasic
	TestServiceControlCache
	TestServiceControlCheckError
	TestServiceControlCheckRetry
	TestServiceControlCheckTimeout
	TestServiceControlCheckTracesWithRetry
	TestServiceControlCheckWrongServerName
	TestServiceControlCredentialId
	TestServiceControlJwtAuthFail
	TestServiceControlLogHeaders
	TestServiceControlLogJwtPayloads
	TestServiceControlNetworkFailFlagClosed
	TestServiceControlNetworkFailFlagOpen
	TestServiceControlProtocolWithGRPCBackend
	TestServiceControlProtocolWithHTTPBackend
	TestServiceControlQuota
	TestServiceControlQuotaExhausted
	TestServiceControlQuotaRetry
	TestServiceControlQuotaUnavailable
	TestServiceControlReportNetworkFail
	TestServiceControlReportResponseCode
	TestServiceControlReportRetry
	TestServiceControlRequestInDynamicRouting
	TestServiceControlRequestWithAllowCors
	TestServiceControlRequestWithoutAllowCors
	TestServiceControlSkipUsage
	TestServiceControlSkipUsageTraces
	TestSimpleCorsWithBasicPreset
	TestSimpleCorsWithRegexPreset
	TestTranscodingBindings
	TestTranscodingErrors
	TestTranscodingServiceUnavailableError
	// The number of total tests. has to be the last one.
	maxTestNum
)

const (
	portBase uint16 = 20000

	// Maximum number of ports used by non-jwt components.
	portNum uint16 = 6
)

var (
	// Maximum number of ports used by jwt fake servers.
	jwtPortNum = uint16(len(testdata.ProviderConfigs))

	// Ports allocated to Jwt open-id servers
	preAllocatedPorts = map[uint16]bool{
		32024: true,
		32025: true,
		32026: true,
		32027: true,
		55550: true,
	}
)

// Ports stores all used ports and other ids for shared resources
type Ports struct {
	BackendServerPort         uint16
	DynamicRoutingBackendPort uint16
	ListenerPort              uint16
	DiscoveryPort             uint16
	AdminPort                 uint16
	FakeStackdriverPort       uint16
	JwtRangeBase              uint16
}

func allocPortBase(testId uint16) uint16 {

	// The maximum number of ports a single test can use
	maxPortsPerTest := portNum + jwtPortNum

	// Find a range of ports that is not taken
	base := portBase + testId*maxPortsPerTest
	for i := 0; i < 10; i++ {
		if allPortFree(base, maxPortsPerTest) {
			return base
		}
		base += maxTestNum * maxPortsPerTest
	}

	glog.Warningf("test(%v) could not find free ports, continue the test...", testId)
	return base
}

func allPortFree(base uint16, ports uint16) bool {
	for port := base; port < base+ports; port++ {
		if IsPortUsed(port) {
			glog.Infoln("port is used ", port)
			return false
		}
	}
	return true
}

// IsPortUsed checks if a port is used
func IsPortUsed(port uint16) bool {

	// Check if this is pre-allocated and should not be used
	_, ok := preAllocatedPorts[port]
	if ok {
		return true
	}

	// Check if anything is listening on this port
	serverPort := fmt.Sprintf("localhost:%v", port)
	conn, _ := net.DialTimeout("tcp", serverPort, 100*time.Millisecond)

	if conn != nil {
		_ = conn.Close()
		return true
	}
	return false
}

// NewPorts allocate all ports based on test id.
func NewPorts(testId uint16) *Ports {
	base := allocPortBase(testId)
	ports := &Ports{
		BackendServerPort:         base,
		DynamicRoutingBackendPort: base + 1,
		ListenerPort:              base + 2,
		DiscoveryPort:             base + 3,
		AdminPort:                 base + 4,
		FakeStackdriverPort:       base + 5,
		JwtRangeBase:              base + 6,
	}
	glog.Infof(fmt.Sprintf("Ports generated for test(%v) are: %+v", testId, ports))
	return ports
}
