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

package error

const (
	// Prefix code used for error codes.
	PREFIX_CODE = "SSRVFCC"

	// DID Document error codes.
	DIDDOC_INSERT_ERROR                 = "01001"
	DIDDOC_GET_ERROR                    = "01002"
	DIDDOC_PUT_ERROR                    = "01003"
	DIDDOC_PROVIDER_INVALID             = "01004"
	DIDDOC_CONVERT_ERROR                = "01005"
	DIDDOC_SIGNATURE_VERIFICATION_ERROR = "01006"
	DID_KEY_URL_PARSING_ERROR           = "01007"
	DIDDOC_VERSIONID_INVAILD            = "01008"

	// DID Document Status error codes.
	DIDDOC_STATUS_INSERT_ERROR  = "02001"
	DIDDOC_STATUS_GET_ERROR     = "02002"
	DIDDOC_STATUS_PUT_ERROR     = "02003"
	DIDDOC_STATUS_CONVERT_ERROR = "02004"
	DIDDOC_STATUS_INVALID       = "02005"

	// VC Meta error codes.
	VCMETA_INSERT_ERROR   = "03001"
	VCMETA_GET_ERROR      = "03002"
	VCMETA_PUT_ERROR      = "03003"
	VCMETA_CONVERT_ERROR  = "03004"
	VCMETA_STATUS_INVALID = "03005"
)

// Mapping of error codes to error messages.
var errMsg = map[string]string{
	DIDDOC_INSERT_ERROR:                 "Failed to insert did document",
	DIDDOC_GET_ERROR:                    "Failed to get did document",
	DIDDOC_PUT_ERROR:                    "Failed to update did document",
	DIDDOC_PROVIDER_INVALID:             "Provider is invalid",
	DIDDOC_CONVERT_ERROR:                "Failed to convert json data to did document",
	DIDDOC_SIGNATURE_VERIFICATION_ERROR: "Failed to signature verification",
	DID_KEY_URL_PARSING_ERROR:           "VerificationMethod is invalid",
	DIDDOC_VERSIONID_INVAILD:            "VersionId is inavlid",

	DIDDOC_STATUS_INSERT_ERROR:  "Failed to insert did document status",
	DIDDOC_STATUS_GET_ERROR:     "Failed to get did document status",
	DIDDOC_STATUS_PUT_ERROR:     "Failed to update did document status",
	DIDDOC_STATUS_CONVERT_ERROR: "Failed to convert json data to did document status",
	DIDDOC_STATUS_INVALID:       "Did document status is invalid",

	VCMETA_INSERT_ERROR:   "Failed to insert vc meta",
	VCMETA_GET_ERROR:      "Failed to get vc meta",
	VCMETA_PUT_ERROR:      "Failed to update vc meta",
	VCMETA_CONVERT_ERROR:  "Failed to convert json data to vc meta",
	VCMETA_STATUS_INVALID: "Vc meta status is invalid",
}

type ContractError struct {
	Code    string
	Message string
}

func (c *ContractError) Error() string {
	return "ErrorCode : " + c.Code + ", Message : " + c.Message
}

func GetContractError(code string, err error) *ContractError {

	return &ContractError{
		Code:    PREFIX_CODE + code,
		Message: errMsg[code] + " => " + err.Error(),
	}
}
