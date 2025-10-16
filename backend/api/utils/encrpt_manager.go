package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"


	"encoding/base64"
	"encoding/hex"

)

var key []byte = []byte("01234567890123456789012345678923")

func calculateMD5(input string) string {

	hasher := md5.New()


	hasher.Write([]byte(input))


	hashBytes := hasher.Sum(nil)


	hashString := hex.EncodeToString(hashBytes)

	return hashString
}

func decrypt(ciphertext string) (string, error) {

	ciphertextBytes, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}


	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}


	iv := ciphertextBytes[:aes.BlockSize]
	ciphertextBytes = ciphertextBytes[aes.BlockSize:]


	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(ciphertextBytes, ciphertextBytes)


	padding := int(ciphertextBytes[len(ciphertextBytes)-1])
	decryptedStr := ciphertextBytes[:len(ciphertextBytes)-padding]


	return string(decryptedStr), nil
}
