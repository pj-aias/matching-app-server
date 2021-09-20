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
		t.Fatalf("auth.ValidatePassword() failed; err = %v", got)
	}

	if cmp.Equal(got, wanted) {
		t.Errorf("invalid result: wanted %v, got %v", wanted, got)
	}
}
