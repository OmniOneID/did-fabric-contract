// Copyright 2024 Raonsecure

package data

const (
	TAS                  ROLE_TYPE = "Tas"
	Wallet               ROLE_TYPE = "Wallet"
	Issuer               ROLE_TYPE = "Issuer"
	WalletProvider       ROLE_TYPE = "WalletProvider"
	AppProvider          ROLE_TYPE = "AppProvider"
	ListProvider         ROLE_TYPE = "ListProvider"
	OpProvider           ROLE_TYPE = "OpProvider"
	KycProvider          ROLE_TYPE = "KycProvider"
	NotificationProvider ROLE_TYPE = "NotificationProvider"
	LogProvider          ROLE_TYPE = "LogProvider"
	PortalProvider       ROLE_TYPE = "PortalProvider"
	DelegationProvider   ROLE_TYPE = "DelegationProvider"
	StorageProvider      ROLE_TYPE = "StorageProvider"
	BackupProvider       ROLE_TYPE = "BackupProvider"
	Etc                  ROLE_TYPE = "Etc"
)

type Provider struct {
	Did       string `validate:"required" json:"did"`
	CertVcRef URL    `validate:"required,url" json:"certVcRef"`
}
