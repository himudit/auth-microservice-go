package config

import (
	"crypto/rsa"
	"log"
	"os"

	"github.com/golang-jwt/jwt/v5"
)

var PrivateKey *rsa.PrivateKey
var PublicKey *rsa.PublicKey

func LoadRSAKeys() {
	privateBytes, err := os.ReadFile("keys/private.pem")
	if err != nil {
		log.Fatal("❌ Cannot read private key:", err)
	}

	publicBytes, err := os.ReadFile("keys/public.pem")
	if err != nil {
		log.Fatal("❌ Cannot read public key:", err)
	}

	PrivateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privateBytes)
	if err != nil {
		log.Fatal("❌ Invalid private key:", err)
	}

	PublicKey, err = jwt.ParseRSAPublicKeyFromPEM(publicBytes)
	if err != nil {
		log.Fatal("❌ Invalid public key:", err)
	}

	log.Println("🔐 RSA keys loaded successfully")
}
