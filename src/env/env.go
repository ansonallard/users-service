package env

import (
	"os"
)

func GetPort() string {
	return os.Getenv("PORT")
}
