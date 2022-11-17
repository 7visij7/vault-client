package aes

import (
	"io"
	"fmt"
	"crypto/aes"
	"crypto/rand"
	"encoding/hex"
	"crypto/cipher"
	"smib-vault-client/types"
	"smib-vault-client/pkg/errors"
)

func  Encrypt(stringToEncrypt string) (encryptedString string) {
	key, _ := hex.DecodeString(types.ENCRYPT_KEY)
	plaintext := []byte(stringToEncrypt)

	block, err := aes.NewCipher(key)
	errors.CheckError(err)

	aesGCM, err := cipher.NewGCM(block)
	errors.CheckError(err)

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}
	
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}
 
func Decrypt(encryptedString string) (decryptedString string) {
	key, _ := hex.DecodeString(types.ENCRYPT_KEY)
	enc, _ := hex.DecodeString(encryptedString)

	block, err := aes.NewCipher(key)
	errors.CheckError(err)

	aesGCM, err := cipher.NewGCM(block)
	errors.CheckError(err)

	nonceSize := aesGCM.NonceSize()

	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	errors.CheckError(err)

	return fmt.Sprintf("%s", plaintext)
}