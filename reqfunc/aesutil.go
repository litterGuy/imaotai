package reqfunc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"imaotai/common/errorx"
)

const (
	aesKey = "qbhajinldepmucsonaaaccgypwuvcjaa"
	aesIV  = "2018534749963515"
)

func AesEncrypt(params string) (string, error) {
	key := []byte(aesKey)
	iv := []byte(aesIV)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errorx.NewDefaultError("Failed to create AES cipher:", err)
	}
	plaintext := []byte(params)
	// 补全填充
	paddingLength := aes.BlockSize - len(plaintext)%aes.BlockSize
	padding := []byte{byte(paddingLength)}
	plaintext = append(plaintext, bytes.Repeat(padding, paddingLength)...)
	// 加密
	ciphertext := make([]byte, len(plaintext))
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, plaintext)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func AesDecrypt(params string) (string, error) {
	key := []byte(aesKey)
	iv := []byte(aesIV)
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", errorx.NewDefaultError("Failed to create AES cipher:", err)
	}
	ciphertext, err := base64.StdEncoding.DecodeString(params)
	if err != nil {
		return "", errorx.NewDefaultError("Failed to decode base64 string:", err)
	}
	// 解密
	plaintext := make([]byte, len(ciphertext))
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(plaintext, ciphertext)
	// 去除填充
	paddingLength := int(plaintext[len(plaintext)-1])
	plaintext = plaintext[:len(plaintext)-paddingLength]
	return string(plaintext), nil
}
