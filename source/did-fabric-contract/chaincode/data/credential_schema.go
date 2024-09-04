// Copyright 2024 Raonsecure

package data

type CREDENTIAL_SCHEMA_TYPE string

type CredentialSchema struct {
	Id   URL                    `validate:"url" json:"id"`
	Type CREDENTIAL_SCHEMA_TYPE `validate:"oneof=OsdSchemaCredential" json:"type"`
}
