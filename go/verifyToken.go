package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"os"
	"strings"
)

func VerifyToken(token string) bool {
	if token == "" {
		return false
	}
	temp := strings.Split(token, ".")
	if len(temp) != 3 {
		return false
	}
	unsignedToken := strings.Join(temp[:2], ".")
	log.Println(unsignedToken)

	// Default secret key or pull from environment
	secretKey := os.Getenv("SECRET_KEY")

	// 1. Create HMAC SHA256 generator
	mac := hmac.New(sha256.New, []byte(secretKey))

	// 2. Hash the unsigned token (equivalent to .update())
	mac.Write([]byte(unsignedToken))
	expectedMAC := mac.Sum(nil)

	// 3. Encode to Base64URL without padding (equivalent to .digest('base64url'))
	expectedSignature := base64.RawURLEncoding.EncodeToString(expectedMAC)

	// 4. constant-time comparison (equivalent to timingSafeEqual)
	return hmac.Equal([]byte(expectedSignature), []byte(temp[2]))

}
