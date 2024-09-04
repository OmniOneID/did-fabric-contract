---
puppeteer:
    pdf:
        format: A4
        displayHeaderFooter: true
        landscape: false
        scale: 0.8
        margin:
            top: 1.2cm
            right: 1cm
            bottom: 1cm
            left: 1cm
    image:
        quality: 100
        fullPage: false
---

ContractError
==

- Topic: ContractError
- Author: Kim Jeong-heon, Kim Min-yong
- Date: 2024-08-29
- Version: v1.0.0

| Version          | Date       | Changes                  |
| ---------------- | ---------- | ------------------------ |
| v1.0.0           | 2024-08-29 | Initial version          |

<div style="page-break-after: always;"></div>

# Table of Contents
- [ContractError](#contracterror)
- [Table of Contents](#table-of-contents)
- [Model](#model)
  - [ContractError](#contracterror-1)
    - [Description](#description)
    - [Declaration](#declaration)
    - [Property](#property)
- [Error Code](#error-code)
  - [1. DID Document(01XXX)](#1-did-document01xxx)
  - [2. DID Document Status(02XXX)](#2-did-document-status02xxx)
  - [3. Verifiable Credential Metadata(03XXX)](#3-verifiable-credential-metadata03xxx)

# Model
## ContractError

### Description
```
Error struct for Fabric Contract. It has code and message pair.
Code starts with SSRVFCC.
```

### Declaration
```go
// Declaration in Golang
type ContractError struct {
    Code string
    Message string
}
```

### Property

| Name           | Type       | Description                            | **M/O** | **Note**         |
|----------------|------------|----------------------------------------|---------|------------------|
| code           | String     | Error code. It starts with SSRVFCC     |    M    |                  | 
| message        | String     | Error description                      |    M    |                  | 

<br>

# Error Code
## 1. DID Document(01XXX)

| Error Code   | Error Message                                 | Description                                                                       | Action Required                                      |
|--------------|-----------------------------------------------|-----------------------------------------------------------------------------------|------------------------------------------------------|
| SSRVFCC01001 | Failed to insert did document                 | Cannot insert a DID document into the blockchain state.                           | Verify DID document data validity and ensure there are no key id conflicts. Depend on detail error cases |
| SSRVFCC01002 | Failed to get did document                    | Cannot retrieve the DID document from the ledger.                                 | Ensure that the DID document status exists in the ledger |
| SSRVFCC01003 | Failed to update did document                 | Cannot update the DID document in the blockchain state.                           | Verify DID document data validity |
| SSRVFCC01004 | Provider is invalid                           | The specified provider is not recognized or invalid.                              | Check the provider's key id and role type |
| SSRVFCC01005 | Failed to convert json data to did document   | JSON data could not be converted to a DID document struct                         | Validate the DID document json format  |
| SSRVFCC01006 | Failed to signature verification              | Signature verification for the DID document failed.                               | Verify the signature value and plain text. Depend on detail error cases |
| SSRVFCC01007 | VerificationMethod is invalid                 | The verification method used is invalid or unsupported.                           | Check the verification method is used.          |
| SSRVFCC01008 | VersionId is invalid                          | The provided VersionId is not valid     | Check versionId of current and previous did documents      |


<br>

## 2. DID Document Status(02XXX)

| Error Code   | Error Message                                      | Description                                                                 | Action Required                                      |
|--------------|----------------------------------------------------|-----------------------------------------------------------------------------|------------------------------------------------------|
| SSRVFCC02001 | Failed to insert did document status               | Cannot insert the DID document status into the blockchain state.            | Verify DID document data validity and ensure there are no key id conflicts. Depend on detail error cases |
| SSRVFCC02002 | Failed to get did document status                  | Cannot retrieve the DID document status from the ledger.                    | Ensure that the DID document status exists in the ledger |
| SSRVFCC02003 | Failed to update did document status               | Cannot update the DID document status in the blockchain state.              | Verify the correctness of the status |
| SSRVFCC02004 | Failed to convert json data to did document status | JSON data could not be converted to the DID document status struct.         | Validate the DID document status JSON format |
| SSRVFCC02005 | Did document status is invalid                     | The provided DID document status is invalid or not recognized.              | Verify the DID document status |


<br>

## 3. Verifiable Credential Metadata(03XXX)

| Error Code   | Error Message                          | Description                                                                       | Action Required                                      |
|--------------|----------------------------------------|-----------------------------------------------------------------------------------|------------------------------------------------------|
| SSRVFCC03001 | Failed to insert vc meta               | Cannot insert VC meta into the blockchain state.                                  | Verify VC meta data validity and ensure there are no key id conflicts. Depend on detail error cases |
| SSRVFCC03002 | Failed to get vc meta                  | Cannot retrieve the VC meta from the ledger.                                      | Ensure that the VC meta exists in the ledger |
| SSRVFCC03003 | Failed to update vc meta               | Cannot update the VC meta in the blockchain state.                                | Verify the correctness of the vc meta status |
| SSRVFCC03004 | Failed to convert json data to vc meta | JSON data could not be converted to the VC meta struct.                           | Validate the VC meta JSON format |
| SSRVFCC03005 | VC meta status is invalid              | The provided VC meta status is not valid or recognized.                           | Verify the VC meta status |