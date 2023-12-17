// Package session provides functions to save and load session data (encrypted JWT to interact with protected api routes).
package session

import (
	"crypto/aes"
	vUnreal "dev.hackerman.me/artheon/veverse-shared/unreal"
	"encoding/base64"
	"fmt"
	"games.launch.launcher/config"
	"games.launch.launcher/crypto"
	ll "games.launch.launcher/logger"
	"io"
	"os"
	"path/filepath"
)

func getSessionEncryptionKey() ([]byte, error) {
	encryptionKey, err := base64.StdEncoding.DecodeString(config.SessionEncryptionKey)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to decode session encryption key: %v\n", err))
		return nil, fmt.Errorf("failed to decode session encryption key: %v", err)
	}
	return encryptionKey, nil
}

func getSessionFilePath(app string) (string, error) {
	dir, err := vUnreal.GetProjectSaveDir(app, config.Configuration)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to get app save dir: %v\n", err))
		return "", fmt.Errorf("failed to get app save dir: %v", err)
	}

	return filepath.Join(dir, ".session.bin"), nil
}

func SaveSession(app string, token string) error {
	sessionPath, err := getSessionFilePath(app)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to get session file path: %v\n", err))
		return fmt.Errorf("failed to get session file path: %v", err)
	}

	err = os.MkdirAll(filepath.Dir(sessionPath), os.ModePerm)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to create session directory: %v\n", err))
		return fmt.Errorf("failed to create session directory: %v", err)
	}

	sessionFile, err := os.Create(sessionPath)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to create session file: %v\n", err))
		return fmt.Errorf("failed to create session file: %v", err)
	}
	defer func(sessionFile *os.File) {
		err := sessionFile.Close()
		if err != nil {
			ll.Logger.Error(fmt.Sprintf("failed to close session file: %v\n", err))
		}
	}(sessionFile)

	encryptionKey, err2 := getSessionEncryptionKey()
	if err2 != nil {
		ll.Logger.Error(fmt.Sprintf("failed to decode session encryption key: %v\n", err2))
		return fmt.Errorf("failed to decode session encryption key: %v", err2)
	}

	sessionBytes, err := encryptSession(encryptionKey, token)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to encrypt session: %v\n", err))
		return fmt.Errorf("failed to encrypt session: %v", err)
	}

	sessionFileHeader := []byte{'M', 'S', 'V', 0x00}
	_, err = sessionFile.Write(sessionFileHeader)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to write session file header: %v\n", err))
		return fmt.Errorf("failed to write session file header: %v", err)
	}

	_, err = sessionFile.Write(sessionBytes)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to write session file: %v\n", err))
		return fmt.Errorf("failed to write session file: %v", err)
	}

	return nil
}

func encryptSession(key []byte, source string) ([]byte, error) {
	// AES block size
	blockSize := aes.BlockSize
	// encoded text length
	sourceLen := len(source) + 4
	// total buffer size
	bufferSize := blockSize + sourceLen + (blockSize - sourceLen%blockSize)
	// allocate space for ciphered data
	buffer := make([]byte, bufferSize)

	// write buffer length as a little-endian uint32
	for i := 0; i < 4; i++ {
		buffer[i] = byte(sourceLen >> uint(8*i))
	}

	// copy source string to buffer
	for i, v := range source {
		buffer[i+4] = byte(v)
	}

	encryptAES, err := crypto.EncryptAES(key, buffer)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to encrypt session: %v\n", err))
		return nil, err
	}

	return encryptAES, nil
}

//region Not used yet

//goland:noinspection GoUnusedExportedFunction
func LoadSession(app string) (string, error) {
	sessionPath, err := getSessionFilePath(app)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to get session file path: %v\n", err))
		return "", fmt.Errorf("failed to get session file path: %v", err)
	}

	sessionFile, err := os.Open(sessionPath)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to open session file: %v\n", err))
		return "", fmt.Errorf("failed to open session file: %v", err)
	}
	defer func(sessionFile *os.File) {
		err := sessionFile.Close()
		if err != nil {
			ll.Logger.Error(fmt.Sprintf("failed to close session file: %v\n", err))
		}
	}(sessionFile)

	sessionBytes, err := io.ReadAll(sessionFile)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to read session file: %v\n", err))
		return "", fmt.Errorf("failed to read session file: %v", err)
	}

	encryptionKey, err2 := getSessionEncryptionKey()
	if err2 != nil {
		ll.Logger.Error(fmt.Sprintf("failed to decode session encryption key: %v\n", err2))
		return "", fmt.Errorf("failed to decode session encryption key: %v", err2)
	}

	session, err := crypto.DecryptAES(encryptionKey, sessionBytes[4:])
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to decrypt session: %v\n", err))
		return "", fmt.Errorf("failed to decrypt session: %v", err)
	}

	return session, nil
}

//goland:noinspection GoUnusedFunction
func decryptSession(key []byte, ciphertextBytes []byte) (string, error) {
	decryptedBytes, err := crypto.DecryptAES(key, ciphertextBytes)
	if err != nil {
		ll.Logger.Error(fmt.Sprintf("failed to decrypt session: %v\n", err))
		return "", fmt.Errorf("failed to decrypt session: %v", err)
	}

	// read buffer length as a little-endian uint32
	sourceLen := 0
	for i := 0; i < 4; i++ {
		sourceLen |= int(decryptedBytes[i]) << uint(8*i)
	}

	return decryptedBytes[4:sourceLen], nil
}

//endregion
