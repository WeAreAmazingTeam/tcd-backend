package helper

import (
	"os"
	"strconv"
)

func GetStrEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		return v
	}
	return v
}

func GetIntEnv(key string) int {
	s := GetStrEnv(key)
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

func GetBoolEnv(key string) bool {
	s := GetStrEnv(key)
	v, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}
	return v
}
