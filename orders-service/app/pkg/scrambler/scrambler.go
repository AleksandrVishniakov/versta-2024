package scrambler

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

type Scrambler interface {
	Decryptor
	Encryptor
}

type Decryptor interface {
	Decrypt([]byte) ([]byte, error)
}

type Encryptor interface {
	Encrypt([]byte) ([]byte, error)
}

type AES256 struct {
	key []byte
}

func NewAES256(key []byte) *AES256 {
	return &AES256{key: key}
}

func (e *AES256) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	padding := aes.BlockSize - (len(data) % aes.BlockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	data = append(data, padText...)

	ciphertext := make([]byte, aes.BlockSize+len(data))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], data)

	str := base64.StdEncoding.EncodeToString(ciphertext)

	return []byte(str), nil
}

func (e *AES256) Decrypt(encryptedData []byte) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(string(encryptedData))
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(e.key)
	if err != nil {
		return nil, err
	}

	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	iv := data[:aes.BlockSize]
	data = data[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(data, data)

	padding := int(data[len(data)-1])
	data = data[:len(data)-padding]

	return data, nil
}
