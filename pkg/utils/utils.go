package pkg

import (
	"crypto/md5"
	"encoding/csv"
	"encoding/hex"
	"io"
	"log"

	"golang.org/x/crypto/bcrypt"
)

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
