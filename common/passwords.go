package common

import (
	"crypto/sha512"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

var passwordCache = map[string]bool{}

func GetPasswordCache() map[string]bool {
	return passwordCache
}

func ClearPasswordCache() {
	passwordCache = map[string]bool{}
}

func CreatePasswordHash(password string) (string, error) {
	hashedPassword := GetSHA512Hash(password)
	hashedBytes, err := bcrypt.GenerateFromPassword(
		hashedPassword,
		bcrypt.DefaultCost,
	)

	return string(hashedBytes), err
}

func CheckPassword(input string, bcryptString string) bool {
	inputHashed := GetSHA512Hash(input)

	if isCorrect, ok := passwordCache[string(inputHashed)]; ok {
		return isCorrect
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(bcryptString),
		inputHashed,
	)

	isCorrect := err == nil
	passwordCache[string(inputHashed)] = isCorrect
	return isCorrect
}

func CheckPasswordHashed(inputHashed []byte, bcryptString string) bool {
	if len(inputHashed) != sha512.Size {
		return false
	}

	if isCorrect, ok := passwordCache[string(inputHashed)]; ok {
		return isCorrect
	}

	err := bcrypt.CompareHashAndPassword(
		[]byte(bcryptString),
		inputHashed,
	)

	isCorrect := err == nil
	passwordCache[string(inputHashed)] = isCorrect
	return isCorrect
}

func CheckPasswordHashedHex(inputHex string, bcryptString string) bool {
	inputHashed, err := hex.DecodeString(inputHex)

	if err != nil {
		return false
	}

	return CheckPasswordHashed(inputHashed, bcryptString)
}

func GetSHA512Hash(input string) []byte {
	hash := sha512.New()
	hash.Write([]byte(input))
	hashedBytes := hash.Sum(nil)
	return hashedBytes
}
