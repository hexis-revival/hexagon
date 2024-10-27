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

	// Unpad the decrypted data
	decrypted, err = UnpadPKCS7(decrypted)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

func UnpadPKCS7(data []byte) ([]byte, error) {
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

	// NOTE: Score data seems to have always 16 bytes of padding
	//       at the end of the data. I am not sure if this is
	//       right, but it works.
	data = data[:len(data)-16]

	// Sometimes we get trailing null bytes, for some reason...
	return bytes.TrimRight(data, "\x00"), nil
}
