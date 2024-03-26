package proofOfAudit

/*
#cgo darwin,arm64 LDFLAGS: ./lib/libprover.a
#cgo linux LDFLAGS: ./lib/libprover.a -ldl -lrt -lm

#include "lib/libprover.h"
#include <stdlib.h>
*/
import "C"

import (
	"encoding/json"
	"errors"

	_ "github.com/0xBow-io/base-eas-asp/pkg/proof-of-audit/lib"
)

type StdBuffer struct {
	Data []byte `json:"data"`
}

type StdIO struct {
	Buffer StdBuffer `json:"buffer"`
}

type SP1Proof struct {
	Proof  string `json:"proof"`
	Stdin  StdIO  `json:"stdIn"`
	Stdout StdIO  `json:"stdOut"`
}

/*
 */
func VerifyCommitment(secret, urlPath, nonce, timestamp, payload, signature, commitmentId, publicID string) (p *SP1Proof, err error) {
	// gemerates proof
	if out := C.generate_sp1_proof_ffi(
		C.CString(secret),
		C.CString(urlPath),
		C.CString(nonce),
		C.CString(timestamp),
		C.CString(payload),
		C.CString(signature),
		C.CString(commitmentId),
		C.CString(publicID),
	); out != nil {
		outStr := C.GoString(out)
		if outStr != "invalid proof" && outStr != "invalid commitment" && outStr != "" {
			err := json.Unmarshal([]byte(outStr), &p)
			return p, err
		}
		return nil, errors.New(outStr)
	}
	return nil, errors.New("missing output")
}
