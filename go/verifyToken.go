package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"os"
	"strings"
)

var secretKey []byte

func init() {
	key := os.Getenv("SECRET_KEY")
	if key == "" {
		key = "my_super_secret_benchmark_key"
	}
	secretKey = []byte(key)
}

func VerifyToken(token string) bool {
	if token == "" {
		return false
	}
	temp := strings.Split(token, ".")
	if len(temp) != 3 {
		return false
	}
	unsignedToken := strings.Join(temp[:2], ".")

	// 1. Create HMAC SHA256 generator
	mac := hmac.New(sha256.New, secretKey)

	// 2. Hash the unsigned token (equivalent to .update())
	mac.Write([]byte(unsignedToken))
	expectedMAC := mac.Sum(nil)

	// 3. Encode to Base64URL without padding (equivalent to .digest('base64url'))
	expectedSignature := base64.RawURLEncoding.EncodeToString(expectedMAC)

	// 4. constant-time comparison (equivalent to timingSafeEqual)
	return hmac.Equal([]byte(expectedSignature), []byte(temp[2]))

}
