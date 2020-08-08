package std

import (
	"fmt"
	"os"
	"strconv"
)

func MustGetIntEnv(key string) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("environment variable `%s` does not exist", key))
	}
	res, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Errorf("cannot cast env(`%s`) to int: %s", key, err.Error()))
	}
	return res
}
