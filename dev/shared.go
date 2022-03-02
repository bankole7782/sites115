package main

import (
  "os"
  "path/filepath"
  "strings"
  "math/rand"
  "time"
)


func GetRootPath() (string, error) {
	hd, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	dd := os.Getenv("SNAP_USER_COMMON")
	if strings.HasPrefix(dd, filepath.Join(hd, "snap", "go")) || dd == "" {
		dd = filepath.Join(hd, "sites115_data")
    os.MkdirAll(dd, 0777)
	}

	return dd, nil
}


func UntestedRandomString(length int) string {
  var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
  const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

  b := make([]byte, length)
  for i := range b {
    b[i] = charset[seededRand.Intn(len(charset))]
  }
  return string(b)
}
