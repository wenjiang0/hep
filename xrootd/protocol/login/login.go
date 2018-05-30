// Copyright 2018 The go-hep Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package login contains the structures describing request and response for login request.
// Login request should be issued prior to most of the other
// requests (see http://xrootd.org/doc/dev45/XRdv310.pdf, p.10).
// As part of the response, SecurityInformation may be provided,
// indicating that an auth request is required. SecurityInformation
// defines the available authentication protocols together with some additional parameters.
// See XRootD protocol specification, page 127 for further information
// about the format of the SecurityInformation.
package login // import "go-hep.org/x/hep/xrootd/protocol/login"

import (
	"os"
)

// RequestID is the id of the request, it is sent as part of message.
// See xrootd protocol specification for details: http://xrootd.org/doc/dev45/XRdv310.pdf, 2.3 Client Request Format.
const RequestID uint16 = 3007

// ResponseLength is the length of the Response assuming that SecurityInformation is empty.
const ResponseLength = 16

// Response is a response for the login request, which contains the session id and the security information.
type Response struct {
	SessionID           [16]byte
	SecurityInformation []byte
}

// Request holds the login request parameters.
type Request struct {
	Pid          int32   // Pid is the process number associated with this connection.
	Username     [8]byte // Username is the unauthenticated name of the user to be associated with the connection.
	Reserved     byte    // Reserved is an area reserved for future use.
	Ability      byte    // Ability are the client's extended capabilities. See xrootd protocol specification, p. 56.
	Capabilities byte    // Capabilities are the Client capabilities. It is 4 for v3.1.0 client without async support.
	Role         byte    // Role is the role being assumed for this login: administrator or regular user.
	TokenLength  int32   // TokenLength is the length of the following token.
	Token        []byte  // Token is the token supplied by the previous redirection response, plus optional elements.
}

// Capabilities for v3.1.0 client without async support.
const clientCapabilities byte = 4

// NewRequest forms a Request according to provided parameters.
func NewRequest(username, token string) Request {
	var usernameBytes [8]byte
	copy(usernameBytes[:], username)

	return Request{
		Pid:          int32(os.Getpid()),
		Username:     usernameBytes,
		Capabilities: clientCapabilities,
		TokenLength:  int32(len(token)),
		Token:        []byte(token),
	}
}