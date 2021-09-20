package auth

import (
	"encoding/json"
	"fmt"
	"os/exec"

	msgpack "github.com/vmihailenco/msgpack/v5"
)

const verifierPath = "../aias-verifier"

type VerifyParams struct {
	Message   []byte `msgpack:"message"`
	Signature string `msgpack:"signature"`
	Gpk       string `msgpack:"gpk"`
}

func VerifySignature(data interface{}, signature, gpk string) (bool, error) {
	encoded, err := json.Marshal(data)
	if err != nil {
		return false, fmt.Errorf("failed to encode data: %v", err)
	}

	params := VerifyParams{
		Message:   encoded,
		Signature: signature,
		Gpk:       gpk,
	}

	dataMsg, err := msgpack.Marshal(&params)
	if err != nil {
		return false, fmt.Errorf("failed to marshal data into MessagePack: %v", err)
	}

	cmd := exec.Command(verifierPath, "verify")
	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		return false, fmt.Errorf("failed to get verifier stdin: %v", err)
	}
	defer cmdStdin.Close()

	_, err = cmdStdin.Write(dataMsg)
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
