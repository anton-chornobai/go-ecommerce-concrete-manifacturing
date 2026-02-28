package infra

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"math/big"
)

type VerificationaCodeManager struct {}

func (v *VerificationaCodeManager) GenerateCode() (string, error) {
	number, err := rand.Int(rand.Reader, big.NewInt(90000))
	if err != nil {
		return "", err
	}

	code := number.Int64() + 10000

	return fmt.Sprintf("%05d", code), nil
}

func (v *VerificationaCodeManager) HashVerificationCode(code string) string {
	hash := sha256.Sum256([]byte(code))
	return hex.EncodeToString(hash[:])
}

func (v *VerificationaCodeManager) CompareHashAndCode(storedHash, userCode string) bool {
	userHash := v.HashVerificationCode(userCode)

	return subtle.ConstantTimeCompare([]byte(storedHash), []byte(userHash)) == 1
}
