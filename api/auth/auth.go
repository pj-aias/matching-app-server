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
		err := errors.New("password length must be less than 72 bytes")
		return nil, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hash, nil
}

func ValidatePassword(hash []byte, password string) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(password))
}

func CreateToken(userId int) (string, error) {
	claims := CustomClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
			Issuer:    "matching-app",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS512, claims)
	signedToken, err := token.SignedString(privKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func ValidateToken(tokenString string) (userId int, err error) {
	token, err := jwt.Parse(tokenString, func(tk *jwt.Token) (interface{}, error) {
		if _, ok := tk.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", tk.Header["alg"])
		}

		return PubKey, nil
	})

	if ve, ok := err.(*jwt.ValidationError); !token.Valid && ok {
		reason := "unknown reason"
		if ve.Errors&jwt.ValidationErrorExpired != 0 {
			reason = "token has expired"
		} else if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			reason = "not a token string: " + tokenString
		}
		log.Printf("authentication failed (%v) for token %v", reason, tokenString)
		return 0, &AuthenticationError{reason}
	} else if !ok {
		// not a validation error
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// All claim values seems to be float. So we must convert to int after assertion to float
		userId = int(claims["UserId"].(float64))
		log.Printf("authentication complete for user %v", userId)
	} else {
		return 0, errors.New("failed to get user id from token")
	}

	return userId, nil
}

func readPrivateKey() error {
	privKeyPath := os.Getenv("SECRET_KEY")
	pubKeyPath := os.Getenv("PUBLIC_KEY")
	privKeyFile, err := os.Open(privKeyPath)

	if errors.Is(err, os.ErrNotExist) {
		// if not generated key pair yet
		privKeyFile, err := os.Create(privKeyPath)
		if err != nil {
			return fmt.Errorf("failed to create private key file: %v", err)
		}
		pubKeyFile, err := os.Create(pubKeyPath)
		if err != nil {
			return fmt.Errorf("failed to create public key file: %v", err)
		}

		privKey, PubKey, err = generateKeyPair(privKeyFile, pubKeyFile)

		return err
	} else if err != nil {
		// normal error
		return fmt.Errorf("failed to read privkey: %v", err)
	}

	// if opened successfuly

	privKeyBuf, err := ioutil.ReadAll(privKeyFile)
	if err != nil {
		return err
	}

	privKey, err = jwt.ParseRSAPrivateKeyFromPEM(privKeyBuf)
	if err != nil {
		return fmt.Errorf("failed to parse RSA private key: %v", err)
	}

	PubKey = &privKey.PublicKey

	return nil
}

func generateKeyPair(privKeyFile, pubKeyFile *os.File) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

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
		Type:  "RSA PUBLIC KEY",
		Bytes: pubKeyBytes,
	}
	err = pem.Encode(pubKeyFile, pubKeyBlock)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, publicKey, nil

}

func init() {
	err := readPrivateKey()
	if err != nil {
		panic(err)
	}
}
