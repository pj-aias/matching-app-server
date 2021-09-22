package auth

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

const verifierPath = "../aias-verifier"

type VerifyParams struct {
	Message   []byte `msgpack:"message"`
	Signature string `msgpack:"signature"`
	Gpk       string `msgpack:"gpk"`
}

type Message = interface{}
type Signature = string
type Gpk = string

func VerifySignature(message Message, signature Signature, gpk Gpk) (bool, error) {
	params, err := fromData(message, signature, gpk)
	if err != nil {
		return false, fmt.Errorf("failed to format message")
	}

	encoded := params.encode()
	if err != nil {
		return false, fmt.Errorf("failed to encode data")
	}

	cmd := exec.Command(verifierPath, "verify")
	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		return false, fmt.Errorf("failed to get verifier stdin: %v", err)
	}
	defer cmdStdin.Close()

	_, err = cmdStdin.Write(encoded)
	if err != nil {
		return false, fmt.Errorf("failed to write verifier stdin: %v", err)
	}

	outBytes, err := cmd.Output()
	out := string(outBytes)
	if out == "OK" {
		// passed verification
		return true, nil
	} else if out == "NG" {
		// failed verification
		return false, nil
	} else if err != nil {
		// some error
		return false, fmt.Errorf("failed to run verifier: %v", err)
	} else {
		return false, fmt.Errorf("unreachable")
	}
}

func fromData(message Message, signature Signature, gpk Gpk) (VerifyParams, error) {
	encoded, err := json.Marshal(message)
	if err != nil {
		return VerifyParams{}, err
	}

	return VerifyParams{
		Message:   encoded,
		Signature: signature,
		Gpk:       gpk,
	}, nil

}

func (p VerifyParams) encode() []byte {
	return append(append([]byte(p.Signature), []byte(p.Gpk)...), p.Message...)
}
