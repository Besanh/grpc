package dotenv

import (
	"os"
	"strconv"
	"time"
)

func GetStringENV(envVar string, defaultValue string) (value string) {
	value = os.Getenv(envVar)
	if len(value) < 1 {
		value = defaultValue
	}
	return
}

func GetIntENV(envVar string, defaultValue int) (value int) {
	value = defaultValue
	if valueStr := os.Getenv(envVar); len(valueStr) > 0 {
		value, _ = strconv.Atoi(valueStr)
	}
	return
}

func GetBoolENV(envVar string, defaultValue bool) (value bool) {
	value = defaultValue
	if valueStr := os.Getenv(envVar); len(valueStr) > 0 {
		value, _ = strconv.ParseBool(valueStr)
	}
	return
}

func GetTimeDurationENV(envVar string, defaultValue time.Duration) (value time.Duration) {
	value = defaultValue
	if valueStr := os.Getenv(envVar); len(valueStr) > 0 {
		value, _ = time.ParseDuration(valueStr)
	}
	return
}
