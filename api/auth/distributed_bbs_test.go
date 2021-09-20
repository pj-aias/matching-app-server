package auth

import (
	"encoding/base64"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFromData(t *testing.T) {
	message := struct {
		hoge string
		fuga int
	}{
		hoge: "hogehoge",
		fuga: 1234,
	}
	signature := "SomeSignature"
	gpk := "SomeGPK"

	got, err := fromData(message, signature, gpk)
	wanted := VerifyParams{
		Message:   []byte(`{"hoge":"hogehoge","fuga":1234}`),
		Signature: signature,
		Gpk:       gpk,
	}

	if err != nil {
		t.Fatalf("auth.fromData() failed; err = %v", err)
	}

	if cmp.Equal(got, wanted) {
		t.Errorf("invalid result: wanted %v, got %v", wanted, got)
	}
}

func TestEncode(t *testing.T) {
	// generated from Rust
	wanted, err := base64.StdEncoding.DecodeString("k5tzb21lTWVzc2FnZa1zb21lU2lnbmF0dXJlp3NvbWVHcGs=")
	if err != nil {
		t.Fatalf("failed to decode sample base64")
	}

	params := VerifyParams{
		Message:   []byte("someMessage"),
		Signature: "someSignature",
		Gpk:       "someGpk",
	}
	got, err := params.encode()

	if err != nil {
		t.Fatalf("auth.VerifyParams.encode() failed; err = %v", err)
	}

	if !cmp.Equal(got, wanted) {
		t.Errorf("invalid result: wanted %v, got %v", wanted, got)
	}
}
