// Copyright 2024 OmniOne.
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

package data

type DID_SERVICE_TYPE string
type DID_SERVICE_ID string

const (
	LinkedDomains      DID_SERVICE_TYPE = "LinkedDomains"
	Credentialregistry DID_SERVICE_TYPE = "Credentialregistry"
)

type Service struct {
	Id              DID_SERVICE_ID   `validate:"required" json:"id"`
	Type            DID_SERVICE_TYPE `validate:"required,oneof=Credentialregistry LinkedDomains" json:"type"`
	ServiceEndpoint []URL            `validate:"required" json:"serviceEndpoint"`
}

func (a *Service) IsEqual(b *Service) bool {
	if a.Id != b.Id {
		return false
	}

	if a.Type != b.Type {
		return false
	}

	if len(a.ServiceEndpoint) != len(b.ServiceEndpoint) {
		return false
	}

	var baseServiceEndpointLength = len(a.ServiceEndpoint)
	for i := 0; i <= baseServiceEndpointLength; i++ {
		if a.ServiceEndpoint[i] != b.ServiceEndpoint[i] {
			return false
		}
	}

	return true
}
