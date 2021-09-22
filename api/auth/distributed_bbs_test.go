package auth

import (
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
	wanted := []byte(`someSignature
someGpk
someMessage`)

	params := VerifyParams{
		Message:   []byte("someMessage"),
		Signature: "someSignature",
		Gpk:       "someGpk",
	}
	got := params.encode()

	if !cmp.Equal(got, wanted) {
		t.Errorf("invalid result: wanted %v, got %v", wanted, got)
	}
}
