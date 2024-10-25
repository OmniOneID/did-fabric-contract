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

import (
	"fmt"
)

type DIDDOC_STATUS string
type ROLE_TYPE string

const (
	DOC_ACTIVATED   DIDDOC_STATUS = "ACTIVATED"
	DOC_DEACTIVATED DIDDOC_STATUS = "DEACTIVATED"
	DOC_REVOKED     DIDDOC_STATUS = "REVOKED"
	DOC_TERMINATED  DIDDOC_STATUS = "TERMINATED"

	DIDDOC_STATUS_PREFIX = "open:did:status:"
)

type DocumentStatus struct {
	Id             string        `json:"did"`
	Status         DIDDOC_STATUS `json:"status"`
	Version        string        `json:"version"`
	Type           ROLE_TYPE     `json:"type"`
	TerminatedTime UTCDateTime   `json:"cancelled_time,omitempty"`
}

func (s *DocumentStatus) Key() ([]string, error) {
	return []string{DIDDOC_STATUS_PREFIX, s.Id}, nil
}

func (s *DocumentStatus) Revoke() error {
	if s.Status == DOC_ACTIVATED || s.Status == DOC_DEACTIVATED {
		s.Status = DOC_REVOKED
		return nil
	}
	return fmt.Errorf("cannot update status from %s to %s", s.Status, DOC_REVOKED)
}

func (s *DocumentStatus) Terminated(terminatedTime UTCDateTime) error {
	if s.Status == DOC_REVOKED {
		s.Status = DOC_TERMINATED
		s.TerminatedTime = terminatedTime
		return nil
	}
	return fmt.Errorf("cannot update status from %s to %s", s.Status, DOC_TERMINATED)
}

func MakeDocumentStatus(didDoc *DidDoc, roleType ROLE_TYPE) *DocumentStatus {
	return &DocumentStatus{
		Id:      didDoc.Id,
		Status:  DOC_ACTIVATED,
		Version: didDoc.VersionId,
		Type:    roleType,
	}
}
