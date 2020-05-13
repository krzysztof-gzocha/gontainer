package std

import (
	"os"
)

func GetEnv(n string) {
	os.Getenv(n)
}
