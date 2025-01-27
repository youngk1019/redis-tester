package internal

import (
	"fmt"
	"os"
	"strings"

	testerutils_random "github.com/codecrafters-io/tester-utils/random"
	"github.com/hdt3213/rdb/parser"
)

func convertToStringArray(interfaceSlice []interface{}) ([]string, error) {
	stringSlice := make([]string, 0, len(interfaceSlice))

	for _, v := range interfaceSlice {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element is not a string: %v", v)
		}
		stringSlice = append(stringSlice, str)
	}

	return stringSlice, nil
}

func compareStringSlices(actual, expected []string, caseSensitiveMatch bool) error {
	if len(actual) != len(expected) {
		return fmt.Errorf("Length mismatch between expected and received messages.\nExpected %v bytes, Received %v bytes.\nExpected : %v.\nReceived : %v.\n", GetByteOffset(expected), GetByteOffset(actual), expected, actual)
	}

	for i := range actual {
		// Wildcard for comparison.
		if expected[i] == "*" {
			continue
		}
		var a, e string
		if caseSensitiveMatch {
			a, e = actual[i], expected[i]
		} else {
			// Case Insensitive matching
			a, e = strings.ToUpper(actual[i]), strings.ToUpper(expected[i])
		}
		if i == 0 {
			// First element in the array is the REDIS command
			// That should always be comapred in a case insensitive manner
			a, e = strings.ToUpper(actual[i]), strings.ToUpper(expected[i])
		}
		if a != e {
			return fmt.Errorf("Expected: '%v' and actual: '%v' messages don't match", e, a)
		}
	}

	return nil
}

func compareStringSlicesWithOr(actual []string, expected [][]string, caseSensitiveMatch bool) error {
	var foundMatch bool
	var e error

	for _, exp := range expected {
		e := compareStringSlices(actual, exp, caseSensitiveMatch)
		if e == nil {
			foundMatch = true
		}
	}

	if foundMatch {
		return nil
	}
	return e // Will return last error. Will accordingly call assert.
}

func parseInfoOutput(lines []string, seperator string) map[string]string {
	infoMap := make(map[string]string)
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		parts := strings.Split(trimmedLine, seperator)
		if len(parts) == 2 {
			key := parts[0]
			value := parts[1]
			infoMap[key] = value
		}
	}
	return infoMap
}

func deleteRDBfile() {
	fileName := "dump.rdb"
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return
	}
	_ = os.Remove(fileName)
}

// Used for parsing RDB file, to check validity.
func processRedisObject(o parser.RedisObject) bool {
	switch o.GetType() {
	case parser.StringType:
		str := o.(*parser.StringObject)
		println(str.Key, str.Value)
	case parser.ListType:
		list := o.(*parser.ListObject)
		println(list.Key, list.Values)
	case parser.HashType:
		hash := o.(*parser.HashObject)
		println(hash.Key, hash.Hash)
	case parser.ZSetType:
		zset := o.(*parser.ZSetObject)
		println(zset.Key, zset.Entries)
	}
	return true
}

func RandomAlphanumericString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		charIndex := testerutils_random.RandomInt(0, len(charset))
		result[i] = charset[charIndex]
	}
	return string(result)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
