package utils

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	uuidLock *sync.Mutex
	lastNum  int64
	count    int
)

func init() {
	uuidLock = new(sync.Mutex)
	count = 0
}

// UUID generate unique values
func UUID() string {
	uuidLock.Lock()
	result := time.Now().UnixNano()
	if lastNum == result {
		count++
	} else {
		count = 0
		lastNum = result
	}
	uuidLock.Unlock()
	return MD5String(strconv.Itoa(int(lastNum)) + strconv.Itoa(count))
}

// MD5String MD5 string
func MD5String(input string) string {
	h := md5.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
}

// JSONStrToSlice Json str to slice
func JSONStrToSlice(jsonStr string) []string {
	jsonStr = strings.Replace(jsonStr, " ", "", -1)
	return strings.Split(strings.Trim(jsonStr, ","), ",")
}
