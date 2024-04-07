package util

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const MAX_LIMIT = 50_000

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")

func ParseLimit(limit any) int {
	var val = 10
	var err error
	if reflect.TypeOf(limit).Kind() == reflect.String {
		str := ParseString(limit)
		val, err = strconv.Atoi(str)
		if err != nil {
			val = 10
		}
	} else {
		val = ParseInt(fmt.Sprintf("%d", limit))
	}
	if val < 0 {
		val = 10
	} else if val > MAX_LIMIT {
		val = MAX_LIMIT
	}
	return val
}

func ParseOffset(offset any) int {
	var val = 0
	var err error
	if reflect.TypeOf(offset).Kind() == reflect.String {
		str := ParseString(offset)
		val, err = strconv.Atoi(str)
		if err != nil {
			val = 0
		}
	} else {
		val = ParseInt(fmt.Sprintf("%d", offset))
	}
	if val < 0 {
		val = 0
	}
	return val
}

func ParseInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return i
}

func ParseAnyToString(value any) (string, error) {
	bytes, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func ParseStringToAny(value string, dest any) error {
	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return err
	}
	return nil
}

func ParseAnyToAny(value any, dest any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(bytes, dest); err != nil {
		return err
	}
	return nil
}

func ParseString(value any) string {
	str, ok := value.(string)
	if !ok {
		return str
	}
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Trim(str, "\r\n")
	str = strings.TrimSpace(str)
	return str
}

func InArray(item any, array any) bool {
	arr := reflect.ValueOf(array)
	if arr.Kind() != reflect.Slice {
		return false
	}
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}

func GenerateRandomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func ParseToAnyArray(value []string) []any {
	result := make([]any, 0)
	for _, v := range value {
		result = append(result, v)
	}
	return result
}

func TimeToString(valueTime time.Time) string {
	return TimeToStringLayout(valueTime, "2006-01-02 15:04:05")
}

func TimeToStringLayout(valueTime time.Time, layout string) string {
	return valueTime.Format(layout)
}

func InArrayContains(item string, array []string) bool {
	for _, v := range array {
		if strings.Contains(item, v) {
			return true
		}
	}
	return false
}

func ParseStructToByte(data any) ([]byte, error) {
	value, err := ParseAnyToString(data)
	if err != nil {
		return nil, err
	}
	byteSlice := []byte(value)
	return byteSlice, nil
}

func ParseQueryArray(slice []string) []string {
	result := make([]string, 0)
	for _, v := range slice {
		if len(v) > 0 {
			result = append(result, v)
		}
	}
	return result
}
