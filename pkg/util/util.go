// SPDX-FileCopyrightText: 2020-present Intel
//
// SPDX-License-Identifier: Apache-2.0
//

package util

import (
	"github.com/free5gc/openapi/AsSessionWithQoS"
)

func GetAsSessionWithQoSClient(uri string) *AsSessionWithQoS.APIClient {
	configuration := AsSessionWithQoS.NewConfiguration()
	configuration.SetBasePath(uri)
	client := AsSessionWithQoS.NewAPIClient(configuration)
	return client
}
