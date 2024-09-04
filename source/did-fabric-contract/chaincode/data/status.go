// Copyright 2024 Raonsecure

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
