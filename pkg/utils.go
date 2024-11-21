package pkg

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Env Helper functions
func GetEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return defaultValue
}

func GetEnvAsIntOrDefault(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func GetEnvAsBoolOrDefault(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		if boolValue, err := strconv.ParseBool(value); err != nil {
			return boolValue
		}
	}

	return defaultValue
}

func GetEnvAsSliceOrDefault(key string, defaultValue []string) []string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return Split(value, ",")
	}
	return defaultValue
}

func Split(s string, sep string) []string {
	if s == "" {
		return []string{}
	}
	return strings.Split(s, sep)
}

func Map[T, U any](arr []T, f func(T) U) []U {
	out := make([]U, len(arr))
	for i := range arr {
		out[i] = f(arr[i])
	}
	return out
}

func CSVToMap(reader io.Reader) []map[string]interface{} {
	r := csv.NewReader(reader)
	var rows []map[string]interface{}
	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if header == nil {
			header = record
		} else {
			dict := map[string]interface{}{}
			for i := range header {
				dict[header[i]] = record[i]
			}
			rows = append(rows, dict)
		}
	}
	return rows
}

func HashPassword(password string) (error, string) {
	// Hash the password
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return err, string(bytes)
}

func ComparePassword(hashedPassword string, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func GenerateMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

func CalculatePercentageChange(from float64, to float64) float64 {
	return ((from - to) / to) * 100
}

func ConvertToUnixTime(timeStr string) time.Time {
	t, _ := strconv.ParseInt(timeStr, 10, 64)
	return time.Unix(t, 0)
}
