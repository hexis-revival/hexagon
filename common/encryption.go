package common

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"errors"
)

func DecryptScoreData(iv []byte, encryptedData []byte) ([]byte, error) {
	if len(encryptedData) == 0 {
		return []byte{}, nil
	}

	if len(iv) != 16 {
		return nil, errors.New("IV must be 16 bytes")
	}

	encryptionKey := "9viq4mujm86947ujxs7i5z82sa6rrzhz"
	data, err := AESDecrypt(encryptionKey, iv, encryptedData)

	if err != nil {
		return nil, err
	}

	return data, nil
}

func AESDecrypt(key string, iv []byte, encryptedData []byte) (decrypted []byte, err error) {
	defer HandlePanic(&err)

	// Ensure the encryption key length is valid for AES-256
	if len(key) != 32 {
		return nil, errors.New("encryption key must be 32 bytes for AES-256")
	}

	if len(iv) != 16 {
		return nil, errors.New("IV must be 16 bytes")
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		return nil, err
	}

	// Create a new CBC decrypter
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the data
	decrypted = make([]byte, len(encryptedData))
	mode.CryptBlocks(decrypted, encryptedData)

	// Hexis base64-encodes an output buffer that contains the encrypted
	// ciphertext plus one extra zero'd 16-byte block at the end.
	// We'll have to remove that padding after decryption.
	decrypted, err = UnpadScoreData(decrypted)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func UnpadScoreData(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("data is empty")
	}

	// Block size for PKCS#7 (128 bits / 8 bits per byte)
	blockSize := 16
	remaining := len(data) % blockSize

	// Check if the data length is a multiple of the block size
	if remaining != 0 {
		return nil, errors.New("data length is not a multiple of the block size")
	}

	if len(data) < blockSize {
		return nil, errors.New("data is too short")
	}

	// The client always appends one extra zero-filled ciphertext
	// block to the serialized encrypted payload
	data = data[:len(data)-blockSize]

	return bytes.TrimRight(data, "\x00"), nil
}
