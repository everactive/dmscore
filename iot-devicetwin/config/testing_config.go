// -*- Mode: Go; indent-tabs-mode: t -*-

/*
 * This file is part of the IoT Identity Service
 * Copyright 2019 Canonical Ltd.
 *
 * This program is free software: you can redistribute it and/or modify it
 * under the terms of the GNU Affero General Public License version 3, as
 * published by the Free Software Foundation.
 *
 * This program is distributed in the hope that it will be useful, but WITHOUT
 * ANY WARRANTY; without even the implied warranties of MERCHANTABILITY,
 * SATISFACTORY QUALITY, or FITNESS FOR A PARTICULAR PURPOSE.
 * See the GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 */

// Package config is configuration for testing
package config

const testCA = `-----BEGIN CERTIFICATE-----
MIIDjDCCAnSgAwIBAgIJAPRcvEcoawtMMA0GCSqGSIb3DQEBCwUAMFsxCzAJBgNV
BAYTAlVLMRMwEQYDVQQIDApTb21lLVN0YXRlMQ8wDQYDVQQHDAZMb25kb24xFzAV
BgNVBAoMDklvVCBNYW5hZ2VtZW50MQ0wCwYDVQQDDARtcXR0MB4XDTE5MDQxMTIw
MzMyM1oXDTI0MDQxMDIwMzMyM1owWzELMAkGA1UEBhMCVUsxEzARBgNVBAgMClNv
bWUtU3RhdGUxDzANBgNVBAcMBkxvbmRvbjEXMBUGA1UECgwOSW9UIE1hbmFnZW1l
bnQxDTALBgNVBAMMBG1xdHQwggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIB
AQC9CHI5kmM1HkpsfBMHH9k+G3432cIpIhP+6AbfL7iKp6BpDZK+U/XWo3SOgaOU
25RCjvQYPjPDDGx2WHT37SPGNOE/kztxUACA8ADk2tt5CveEfEfHiasNz4ojb2vo
eMPt3q9TdY4PjfzH5Q4kLHMGmY1FXka1dF+WSaB1Psjm/UEqoYE3filtJND3E316
y2Js3Nk3WVTu+ke91jgZM4GV0xEnKlE+MdfN3odaMT6xrUVlDXaQiPzxznE/PNZW
mMbLyKFB/HPWKguPzUvTH4aFcCh+WRF/Am1oPJXfGzPLbzlO81+Wet6RgO0ZSX4E
ZpGiqQ1oAmKfrf6nYzqt7qJBAgMBAAGjUzBRMB0GA1UdDgQWBBRfjuFAbwsbM7l6
Y+IanaIccY2NETAfBgNVHSMEGDAWgBRfjuFAbwsbM7l6Y+IanaIccY2NETAPBgNV
HRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQBv47iHMn9SSfykhOdSmHbx
Oh5Ubw/3jU5WSdsC87CpYUpEJCJUPwz+/d7BbmvcLLnQ62j/gUKX9FgHGGKmpPS0
az/xdg+CG6/ICVf7GtFlt9pXldkUknpV/1l/kN6fI2BPyveB+82ECCvM77RRBh5f
8/342/lw4W6mwHz2esuot4VU7rBRhrJcxAi8xEIjcJx+VJSpQ6p8pM/Pxh2icrWO
PoDcmqYTOpA4A3nyYwJiA8Ph4tO01sOfqTw5geyst4+s4hs1zeoBhdExnWC7s+ti
3wDt/QYhvisV61YYdtmcdFxg44+5Aq6Fc5oiARBt1w+VTo5GxFsPlldY/Vjh8z91
-----END CERTIFICATE-----`

const testServerCert = `-----BEGIN CERTIFICATE-----
MIIDLzCCAhcCCQDSj1KyVrGW6jANBgkqhkiG9w0BAQsFADBbMQswCQYDVQQGEwJV
SzETMBEGA1UECAwKU29tZS1TdGF0ZTEPMA0GA1UEBwwGTG9uZG9uMRcwFQYDVQQK
DA5Jb1QgTWFuYWdlbWVudDENMAsGA1UEAwwEbXF0dDAeFw0xOTA0MTEyMDM2MzJa
Fw0yMDA0MDUyMDM2MzJaMFgxCzAJBgNVBAYTAlVLMRMwEQYDVQQIDApTb21lLVN0
YXRlMQ8wDQYDVQQHDAZMb25kb24xFDASBgNVBAoMC1NlcnZlciBDZXJ0MQ0wCwYD
VQQDDARtcXR0MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvmVaIIuF
ssJEiHTTCX4IXRyYhGU3qacLoNzKoPhTX+uWPpY8tYwpHlyCfgEAguOSbayVfuho
0dyUqaf1bByhMrGxYpAej0i5gZHWzl3eB27u5d1fAkfDpnO6ZIjU2wW3R+rmMsxn
ftUUkJRK06dTqew9q1R4FJH9zS0mjX6GX0p6YqyZrTLbeiNR31Grz9jzUcfZYQit
r/gWqEhHZP59eltUp7bOIvpdRAacqkfCUV8903XShrQM3khp/Hg1qqEY5SlmrvfX
a8RMTBeV6qA18MYw5jXoeMGWqy/u+tcss2AdAGJwv9tKLVSHc+qKnZOi6uNE9GPN
COT72JqF71m6WQIDAQABMA0GCSqGSIb3DQEBCwUAA4IBAQCRynC39KXeie23eQS7
jh0GJATB8U3rtYXfIWk3sVtW0YnkcPMAtqpzLPUrjLwxICFsAMp1+lzLCI/fV5+t
kLrAIFOX22tFlBTiMTCBuQeM8YphAcwijzJSPhJ7MkqyUvPvtvW4TxRDZnu9Txp7
DEV/QqPg58pfRFe5iVTJNSZIof9+tJuAFbVAxTanWrejT+jMplm29HkL4AdhaydZ
8sEhbmrOLdrjYWU78Ytvu7L4DTCOAeRT3vJtCuj+D2p53jHUdqiiSvNoPOcgTR+1
OAryHgJKL9U3TnBgGsm674wx7O6ZukO+x3D1LKGNbmsJ0rdq/Mj8JBdhP5pDoBG6
Jou0
-----END CERTIFICATE-----`

const testServerKey = `-----BEGIN RSA PRIVATE KEY-----
MIIEpQIBAAKCAQEAvmVaIIuFssJEiHTTCX4IXRyYhGU3qacLoNzKoPhTX+uWPpY8
tYwpHlyCfgEAguOSbayVfuho0dyUqaf1bByhMrGxYpAej0i5gZHWzl3eB27u5d1f
AkfDpnO6ZIjU2wW3R+rmMsxnftUUkJRK06dTqew9q1R4FJH9zS0mjX6GX0p6YqyZ
rTLbeiNR31Grz9jzUcfZYQitr/gWqEhHZP59eltUp7bOIvpdRAacqkfCUV8903XS
hrQM3khp/Hg1qqEY5SlmrvfXa8RMTBeV6qA18MYw5jXoeMGWqy/u+tcss2AdAGJw
v9tKLVSHc+qKnZOi6uNE9GPNCOT72JqF71m6WQIDAQABAoIBAQCy94wPWXbUQB2x
crbIjmqIM4/9qzL2Sqn4jHH/e0zLtiQlMo1gTZ59BpI2pPR5FDcY1ogzoXyd/8zR
6Kod9I9lmnfV4QiIwOB2tcKHet5weEshUMO03gY/mTrUs3X5ZtcQR/IYP+Ds7JgH
Cw2HBBr1d7XELYMuOsiqK0245Pyj5jrALhkf9H+qEFTb2GnwKZWlOxey4rGym1Zv
ivxZzRQ41hzyi0Xf6yrDfda87Ns5514M763g4Vfsbb/RKTQLgOfFZJ+oFTuzj2Gz
ezzo6yTfbPUwdlClRnXa/Lvf9xN0qj0dzpKfQh9pHxVZzEAmkNkv4nfE0qT+wwmP
bKT5f9wBAoGBAPARmBNdhvTnTrvSEoTi23a+T22O2VBe1b/ycfOW2NcEI7qWxDXV
vps058YtZREnckE/avpjhaKUbqNOpFP3uOTFjNxD7vGsSFsrk2vjBz8ByrGEPg6H
oYBTOn6L+484jw1CFis1p4YLlN/cBrAQg9YJGlfl4Rxxoz+yQxDNZ1OBAoGBAMsH
5EJ/zb7tKE9ntJ/oE3dz8mCIzqXZJv0bWof/Mb9JtolrjVuLoEWZWNmCjSS0PmmX
i/PZXRzqYE0vXxQCHjFL4UEH206KqmSAK8wHRe1rPJ9L2Hi9Leu+A8OmqPf2LIph
W7mFNlxSqJOQSXFX6TubIuBVxh2oe1XAuh5/1PLZAoGBAOKvClVG3BdGjs7FJx15
hNeUDjYaS9MbKWStDrJ/PtORIhefIzjeUrQFedFkrelLwRQhSOeTr+z7kZj8uihb
YqgKbd7S+r4S+uOzuumFnyL8kyOaBmr74SDl9fbmQSxUsKdJPtugN0ZYi0PyZBI+
Fe61+70B4NVV7FtJ/Q/RlH6BAoGAL3LEsZXUq44ZIZWG7Of7xKrgNhdC1BePuQ8v
dSD6q026zxrHimFzL1DLJuoPukg1XdAA8RgXXq6XmvI7Mh3cmIC3P89qPUzCzYH2
ulPoz7eED2ZWTMFJfhKGJq9IRcrOVfiyywSK08CtjO7newmkhD2ZRPxGtJ+vUzcb
SA1v4uECgYEA5RHoR6LiB84uGrRBvYdugzD/hmjZpBpcchH3J4FSkgjw3UAlrnaI
Sq3Ei5frAGCkh+sTllLtYK2i2EFq67oeytyURH2LM84tGHqZeTdmj/ZKTGuVKW0n
Qmery+icg5rZtFUNYs0WDeeYRVic+fwpBXCjrSpG3KglcceN6uOZVFI=
-----END RSA PRIVATE KEY-----`

// TestMQTTConnect creates config settings for testing
func TestMQTTConnect() *MQTTConnect {
	return &MQTTConnect{
		ClientID:   "aaa",
		RootCA:     []byte(testCA),
		ClientCert: []byte(testServerCert),
		ClientKey:  []byte(testServerKey),
	}
}
