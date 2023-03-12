package constant

import "os"

var SecretKey []byte

func InitAuthConstant() {
	SecretKey = []byte(os.Getenv("SECRET_KEY"))
}
