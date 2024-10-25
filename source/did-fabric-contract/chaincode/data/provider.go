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
