package configh

import (
	"os"
	"strconv"
)

func CombineParams[T comparable](defaultValue *T, params ...*T) *T {
	for i := len(params) - 1; i >= 0; i-- {
		if params[i] != nil {
			return params[i]
		}
	}
	return defaultValue
}

func GetEnvString(key string) *string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return nil
	}
	return &value
}

func GetEnvInt(key string) *int {
	value, exists := os.LookupEnv(key)
	if !exists {
		return nil
	}
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return nil
	}
	return &intValue
}

func GetEnvBool(key string) *bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return nil
	}

	boolValue, err := strconv.ParseBool(value)
	if err != nil {
		return nil
	}

	return &boolValue
}
