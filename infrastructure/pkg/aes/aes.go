package aes

import (
	"bytes"
	cryptoAes "crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

var _ Aes = (*aes)(nil)

type Aes interface {
	i()
	// Encrypt 加密
	Encrypt(encryptStr string) (string, error)

	// Decrypt 解密
	Decrypt(decryptStr string) (string, error)
}

type aes struct {
	key string
	iv  string
}

func NewAes(key, iv string) Aes {
	return &aes{
		key: key,
		iv:  iv,
	}
}

func (a *aes) i() {}

func (a *aes) Encrypt(encryptStr string) (string, error) {
	encryptBytes := []byte(encryptStr)
	block, err := cryptoAes.NewCipher([]byte(a.key))
	if err != nil {
		return "", err
	}

	blockSize := block.BlockSize()
	encryptBytes = pkcs5Padding(encryptBytes, blockSize)

	blockMode := cipher.NewCBCEncrypter(block, []byte(a.iv))
	encrypted := make([]byte, len(encryptBytes))
	blockMode.CryptBlocks(encrypted, encryptBytes)
	return base64.StdEncoding.EncodeToString(encrypted), nil
}

func (a *aes) Decrypt(decryptStr string) (string, error) {
	decryptBytes, err := base64.StdEncoding.DecodeString(decryptStr)
	if err != nil {
		return "", err
	}

	block, err := cryptoAes.NewCipher([]byte(a.key))
	if err != nil {
		return "", err
	}

	blockMode := cipher.NewCBCDecrypter(block, []byte(a.iv))
	decrypted := make([]byte, len(decryptBytes))

	blockMode.CryptBlocks(decrypted, decryptBytes)
	decrypted, err = pkcs5UnPadding(decrypted)
	if err != nil {
		return "", err
	}
	return string(decrypted), nil
}

func pkcs5Padding(cipherText []byte, blockSize int) []byte {
	padding := blockSize - len(cipherText)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(cipherText, padText...)
}

func pkcs5UnPadding(decrypted []byte) ([]byte, error) {
	length := len(decrypted)
	unPadding := int(decrypted[length-1])
	if length < unPadding {
		return nil, errors.New("decrypted error, check key or iv")
	}
	return decrypted[:(length - unPadding)], nil
}
