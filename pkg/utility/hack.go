package utility

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Hack string

// Challenge 正確數值依照每個小時變化, 避免被有心人紀錄
func (hack Hack) Challenge(value string) bool {
	return hack.Value() == value
}

func (hack Hack) Value() string {
	const MMDDHH = "010215"
	cipher := time.Now().Format(MMDDHH)
	m, _ := strconv.Atoi(cipher[0:2])
	d, _ := strconv.Atoi(cipher[2:4])
	h, _ := strconv.Atoi(cipher[4:6])

	key := string(hack)
	const AlphabetSize int = 26
	if len(key) < AlphabetSize {
		key = "abcdefghkjilmnopqrstuvwxyz"
	}

	char1 := key[h%m]
	char2 := key[d%(h+1)]
	char3 := key[(d+h)%len(key)]

	hash := sha256.New()
	hash.Write([]byte{char1, char2, char3})
	hashBytes := hash.Sum(nil)
	length := 10
	if len(hashBytes) > length {
		hashBytes = hashBytes[:length]
	}

	return hex.EncodeToString(hashBytes)
}
