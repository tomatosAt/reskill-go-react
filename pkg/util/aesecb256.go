package util

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"

	"github.com/andreburgaud/crypt2go/ecb"
	"github.com/andreburgaud/crypt2go/padding"
)

func EncryptList(key string, data ...string) ([]string, error) {
	for i, v := range data {
		encryptData, err := EncryptAes256Ecb(key, v)
		if err != nil {
			return nil, err
		}
		data[i] = encryptData
	}
	return data, nil
}

func DecryptList(key string, data ...string) ([]string, error) {
	for i, v := range data {
		decryptData, err := DecryptAes256Ecb(key, v)
		if err != nil {
			return nil, err
		}
		data[i] = decryptData
	}
	return data, nil
}

func EncryptAes256Ecb(key string, payload string) (string, error) {
	payloadByte := []byte(payload)
	secretByte := []byte(key)
	enCode, err := encrypt(payloadByte, secretByte)
	if err != nil {
		return "", fmt.Errorf("key_is_not_valid")
	}
	ciphertext := base64.StdEncoding.EncodeToString(enCode)
	return ciphertext, nil
}

func DecryptAes256Ecb(key string, encryptStr string) (string, error) {
	secretByte := []byte(key)
	sDec, err := base64.StdEncoding.DecodeString(encryptStr)
	if err != nil {
		return "", fmt.Errorf("key_is_not_valid")
	}
	plaintextByte, err := decrypt(sDec, secretByte)
	if err != nil {
		return "", fmt.Errorf("key_is_not_valid")
	}
	return string(plaintextByte), nil
}

func encrypt(pt, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), fmt.Errorf("key_is_not_valid")
	}
	mode := ecb.NewECBEncrypter(block)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Pad(pt)
	if err != nil {
		return []byte(""), fmt.Errorf("key_is_not_valid")
	}
	ct := make([]byte, len(pt))
	mode.CryptBlocks(ct, pt)
	return ct, nil
}

func decrypt(ct, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return []byte(""), fmt.Errorf("key_is_not_valid")
	}
	mode := ecb.NewECBDecrypter(block)
	pt := make([]byte, len(ct))
	mode.CryptBlocks(pt, ct)
	padder := padding.NewPkcs7Padding(mode.BlockSize())
	pt, err = padder.Unpad(pt)
	if err != nil {
		return []byte(""), fmt.Errorf("key_is_not_valid")
	}
	return pt, nil
}
