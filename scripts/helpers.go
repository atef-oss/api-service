// Package helpers provides common utility functions used across the API service.
package helpers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis/v8"
)

// AuthMiddleware adds authentication to the request context.
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := verifyToken(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "token", token)
		next(w, r.WithContext(ctx))
	}
}

func verifyToken(r *http.Request) (jwt.Token, error) {
	tokenString := r.Header.Get("Authorization")
	if tokenString == "" {
		return nil, errors.New("token is missing")
	}

	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

// RedisClient returns a Redis client instance.
func RedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: strings.TrimSpace(os.Getenv("REDIS_PASSWORD")),
		DB:       0,
	})
}

// GetInt returns the integer value of a string or 0 if it's not an integer.
func GetInt(s string) int {
	v, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return v
}

// GetTime returns the parsed time based on the given input string.
func GetTime(s string) time.Time {
	t, err := time.Parse("2006-01-02T15:04:05Z", s)
	if err != nil {
		log.Println(err)
		return time.Now()
	}
	return t
}

// MarshalJSON is a custom JSON marshaler for the NullString type.
func (ns NullString) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(null)
}

type NullString struct {
	String  string
	Valid   bool
}

func (ns NullString) String() string {
	return ns.String
}

func (ns NullString) MarshalText() ([]byte, error) {
	return []byte(ns.String), nil
}

type null struct {
	String  string
	Valid   bool
}

func (ns null) MarshalJSON() ([]byte, error) {
	if ns.Valid {
		return json.Marshal(ns.String)
	}
	return json.Marshal(null)
}

// NullInt is a custom JSON marshaler for the NullInt type.
func (ni NullInt) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.Int)
	}
	return json.Marshal(null)
}

type NullInt struct {
	Int    int
	Valid  bool
}

func (ni NullInt) String() string {
	return strconv.Itoa(ni.Int)
}

func (ni NullInt) MarshalText() ([]byte, error) {
	return []byte(strconv.Itoa(ni.Int)), nil
}

type null struct {
	String  string
	Valid   bool
}

func (ni null) MarshalJSON() ([]byte, error) {
	if ni.Valid {
		return json.Marshal(ni.String)
	}
	return json.Marshal(null)
}