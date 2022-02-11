package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/malyg1n/shortener/pkg/config"
)

const (
	encryptOpType = "encrypt"
	decryptOpType = "decrypt"
)

// Encrypt message.
func Encrypt(msg string) (string, error) {
	out, err := doOperation(encryptOpType, []byte(msg))
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(out), nil
}

// Decrypt message.
func Decrypt(msg string) (string, error) {
	msgBts, err := hex.DecodeString(msg)
	if err != nil {
		return "", err
	}

	out, err := doOperation(decryptOpType, msgBts)
	return string(out), err
}

func doOperation(opType string, msg []byte) ([]byte, error) {
	cfg := config.GetConfig()
	key := sha256.Sum256([]byte(cfg.SecretKey))

	block, err := aes.NewCipher(key[:])
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}

	nonce := key[len(key)-gcm.NonceSize():]
	switch opType {
	case encryptOpType:
		return gcm.Seal(nil, nonce, msg, nil), nil
	case decryptOpType:
		return gcm.Open(nil, nonce, msg, nil)
	}

	return []byte{}, fmt.Errorf("undefined operation type")
}
