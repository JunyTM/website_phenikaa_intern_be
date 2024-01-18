package utils

import (
	"crypto/rand"
	"encoding/base64"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/lithammer/shortuuid"
)

// GenerateKey using in set keys
func GenCode() string {
	id := shortuuid.New()
	return strings.ToUpper(id[0:10])
}

// PatternGet using in get keys
func PatternGet(id uint) string {
	return strconv.Itoa(int(id)) + "-:--*"
}

func InArray(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1
	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				index = i
				exists = true
				return
			}
		}
	}
	return
}

func isMn(r rune) bool {
	return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
}

// GenerateKey random password
func GeneratePasswordKey(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(buffer)[:length], nil
}
