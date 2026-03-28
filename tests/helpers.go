package api_service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/urfave/slug"
)

type RandomData struct {
	Numbers  []int
	Letters  []string
}

func GenerateRandomData() (*RandomData, error) {
	numbers := make([]int, 10)
	letters := make([]string, 10)

	for i := range numbers {
		num, err := rand.Int(rand.Reader, big.NewInt(100))
		if err != nil {
			return nil, err
		}
		numbers[i] = num.Int64()
	}

	for i := range letters {
		letter := make([]byte, 1)
		_, err := rand.Read(letter)
		if err != nil {
			return nil, err
		}
		letters[i] = string(letter)
	}

	return &RandomData{numbers, letters}, nil
}

func Slugify(s string) string {
	return strings.ToLower(slug.Make(s))
}

func GenerateUUID() string {
	uuid := make([]byte, 16)
	_, err := rand.Read(uuid)
	if err != nil {
		panic(err)
	}
	uuid[6] = (uuid[6] & 0x0f) | 0x40
	uuid[8] = (uuid[8] & 0x3f) | 0x80
	return fmt.Sprintf("%x-%x-%x-%x-%x",
		uuid[0:4], uuid[4:6], uuid[6:8], uuid[8:10], uuid[10:])
}

func ParseUUID(uuid string) (string, error) {
	if !strings.HasPrefix(uuid, "uuid-") {
		return "", fmt.Errorf("uuid must be in uuid-123e4567-e89b-12d3-a456-426655440000 format")
	}
	id := strings.TrimPrefix(uuid, "uuid-")
	if len(id) != 32 {
		return "", fmt.Errorf("uuid must be in uuid-123e4567-e89b-12d3-a456-426655440000 format")
	}
	return id, nil
}

func ParseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func IsStringValid(s string) bool {
	return len(s) > 0 && !strings.ContainsAny(s, "^\n\r")
}

func CopyHeader(src, dst http.Header) {
	for k, v := range src {
		dst[k] = v
	}
}

func GetRandomBytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func IsStringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func EncodeCookie(cookie securecookie.Cookie, value interface{}) ([]byte, error) {
	return cookie.Encrypt([]byte(fmt.Sprintf("%v", value)))
}

func DecodeCookie(cookie securecookie.Cookie, value []byte) (interface{}, error) {
	var data interface{}
	err := cookie.Decode(string(value), &data)
	return data, err
}

func IsTimeInPast(t time.Time) bool {
	return t.Before(time.Now())
}