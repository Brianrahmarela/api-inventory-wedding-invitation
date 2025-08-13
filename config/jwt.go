package config

import (
	"os"
	"time"
)

// ambil JWT Secret keys
func GetJwtSecret() []byte {
	return []byte(os.Getenv("JWT_SECRET_KEY"))
}

// ambil durasi expired dari JWT
func GetJwtExpirationDuration() time.Duration {
	duration, err := time.ParseDuration(os.Getenv("JWT_EXPIRES_IN"))
	// fmt.Println("duration", duration)
	// set default 24 jam jika di .env JWT_EXPIRES_IN == nil
	if err != nil {
		return time.Hour * 24
	}
	return duration
}
