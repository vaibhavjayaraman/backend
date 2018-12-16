package tools

import "os"

var JwtSecretKey = []byte(GetEnv("jwt_secret_key", "aVerySecretKey"))

func GetEnv(key, defaultVal string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultVal
	}
	return key
}