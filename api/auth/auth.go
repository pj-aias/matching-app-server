package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

var (
	privKey *rsa.PrivateKey
	PubKey  *rsa.PublicKey
)

type CustomClaims struct {
	UserId int
	jwt.StandardClaims
}

func GeneratePasswordHash(password string) ([]byte, error) {
	// bcrypt won't work correctly if the password length is > 72
	if len(password) > 72 {
		err := errors.New("Password length must be less than 72 bytes.")
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func Validate(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}

func CreateToken(user_id int) string {
	TODO check
	claims := CustomClaims{
		user_id,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
			Issuer:    "makino",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	signed_string, err := token.SignedString(sign_key)
	if err != nil {
		panic(err)
	}

	return signed_string
}

func VerifyUser(token_value string) (user_id int) {
	TODO Check
	token, err := jwt.Parse(token_value, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return verify_key, nil
	})

	if err != nil {
		panic(err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// All claim values seems to be float. So we must convert to int after assertion to float
		user_id = int(claims["UserId"].(float64))
		log.Printf("authorization complete for user %v", user_id)
	} else {
		panic("Get user_id failed")
	}

	return
}

func readPrivateKey() error {
	privKeyPath := os.Getenv("SECRET_KEY")
	pubKeyPath := os.Getenv("PUBLIC_KEY")
	privKeyFile, err := os.Open(privKeyPath)

	if err == os.ErrNotExist {
		// if not generated key pair yet
		privKeyFile, err = os.Create(privKeyPath)
		if err != nil {
			return err
		}
		pubKeyFile, err := os.Create(pubKeyPath)
		if err != nil {
			return err
		}

		privKey, PubKey, err = generateKeyPair(privKeyFile, pubKeyFile)

		return err
	} else if err != nil {
		// normal error
		return err
	}

	// if opened successfuly

	privKeyBuf, err := ioutil.ReadAll(privKeyFile)
	pubKeyBuf, err := ioutil.ReadAll(pubKeyFile)

	privKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyBuf)
	if err != nil {
		return err
	}

	PubKey, err = jwt.ParseRSAPublicKeyFromPEM(pubKeyBuf)
	if err != nil {
		return err
	}

	return nil
}

func generateKeyPair(privKeyFile, pubKeyFile *os.File) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey

	privKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privKeyBytes,
	}
	err = pem.Encode(privKeyFile, privKeyBlock)
	if err != nil {
		return nil, nil, err
	}

	pubKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	if err != nil {
		return nil, nil, err
	}
	pubKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	err = pem.Encode(pubKeyFile, pubKeyBlock)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil

}

func init() {
	TODO check
	privKeyPath := os.Getenv("SECRET_KEY")
	privKeyFile, err := os.Open(privKeyPath)

	if err == os.ErrNotExist {
		privKeyFile, err = os.Create(privKeyPath)
		rsa.GenerateKey(rand.Reader, 2048)
	}

	if err != nil {

	} else {

	}

	sign_key, err = jwt.ParseRSAPrivateKeyFromPEM(secret_key)
	if err != nil {
		panic(err)
	}

	public_key, err := ioutil.ReadFile(os.Getenv("PUBLIC_KEY_PKCS8"))
	if err != nil {
		panic(err)
	}

	verify_key, err = jwt.ParseRSAPublicKeyFromPEM(public_key)
	if err != nil {
		panic(err)
	}
}
