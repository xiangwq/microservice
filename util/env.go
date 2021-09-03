package util

import (
	"os"
	"strings"
)

const (
	MICROSERVICE_ENV = "MICROSERVICE_ENV"
	PRODUCT_ENV      = "product"
	TEST_ENV         = "test"
)

var (
	cur_microservice_env string = TEST_ENV
)

func init() {
	cur_microservice_env = strings.ToLower(os.Getenv(MICROSERVICE_ENV))
	cur_microservice_env = strings.TrimSpace(cur_microservice_env)

	if len(cur_microservice_env) == 0 {
		cur_microservice_env = TEST_ENV
	}
}

func IsProduct() bool {
	return cur_microservice_env == PRODUCT_ENV
}

func IsTest() bool {
	return cur_microservice_env == TEST_ENV
}

func GetEnv() string {
	return cur_microservice_env
}
