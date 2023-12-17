// Package crypto provides a function to encrypt and decrypt data using AES block cipher.
package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	ll "games.launch.launcher/logger"
	"io"
)

func EncryptAES(key []byte, buffer []byte) ([]byte, error) {

	// validate buffer size
	if len(buffer)%aes.BlockSize != 0 {
		ll.Logger.Error(fmt.Sprintf("plaintext is not a multiple of the block size\n"))
		return nil, fmt.Errorf("plaintext is not a multiple of the block size")
	}

	// create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to create cipher: %v\n", err))
		return nil, fmt.Errorf("failed to create cipher: %v", err)
	}

	// The IV needs to be unique, but not secure.
	// Therefore, it's common to include it at the beginning of the ciphertext.
	ciphertext := make([]byte, aes.BlockSize+len(buffer))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to read random bytes: %v\n", err))
		return nil, fmt.Errorf("failed to read random bytes: %v", err)
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext[aes.BlockSize:], buffer)

	// return encoded bytes
	return ciphertext, nil
}

func DecryptAES(key []byte, ciphertextBytes []byte) (string, error) {
	// create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to create cipher: %v\n", err))
		return "", fmt.Errorf("failed to create cipher: %v", err)
	}

	// allocate space for deciphered data
	buffer := make([]byte, len(ciphertextBytes))

	iv := ciphertextBytes[:aes.BlockSize]
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(buffer, ciphertextBytes)

	return string(buffer[aes.BlockSize:]), nil
}
