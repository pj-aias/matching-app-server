package distributed_bbs

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

const verifierPath = "/usr/local/bin/aias-verifier"

type VerifyParams struct {
	Message   []byte `msgpack:"message"`
	Signature string `msgpack:"signature"`
	Gpk       string `msgpack:"gpk"`
}

type Message = interface{}
type Signature = string
type Gpk = string
type Gm = string

func VerifySignature(message Message, signature Signature, gpk Gpk) (bool, error) {
	params, err := fromData(message, signature, gpk)
	if err != nil {
		return false, fmt.Errorf("failed to format message")
	}

	return params.verify()
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

func (p VerifyParams) verify() (bool, error) {
	encoded := p.encode()

	cmd := exec.Command(verifierPath, "verify")
	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		return false, fmt.Errorf("failed to get verifier stdin: %v", err)
	}

	_, err = cmdStdin.Write(encoded)
	if err != nil {
		return false, fmt.Errorf("failed to write verifier stdin: %v", err)
	}
	cmdStdin.Close()

	outBytes, err := cmd.Output()
	out := string(outBytes)
	if strings.HasPrefix(out, "OK") {
		// passed verification
		return true, nil
	} else if strings.HasPrefix(out, "NG") {
		// failed verification
		return false, nil
	} else if err != nil {
		// some error
		return false, fmt.Errorf("failed to run verifier: %v", err)
	} else {
		return false, fmt.Errorf("unreachable (output: '%v')", out)
	}

}

func (p VerifyParams) encode() []byte {
	/* concat three byte slices
	{SIGNATURE}
	{GPK}
	{MESSAGE}
	*/
	res := []byte(p.Signature)
	res = append(res, '\n')
	res = append(res, []byte(p.Gpk)...)
	res = append(res, '\n')
	res = append(res, p.Message...)
	return res
}
