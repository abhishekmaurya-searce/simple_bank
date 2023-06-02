package util

import (
	"math/rand"
	"time"
)

const alpha = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var str string
	k := len(alpha)
	for i := 0; i < n; i++ {
		str += string(alpha[rand.Intn(k)])
	}
	return str
}
func RandomOwner() string {
	return RandomString(6)
}
func RandomMoney() int64 {
	return RandomInt(0, 10000)
}
func RandomCurrency() string {
	curr := []string{"EUR", "USD", "INR"}
	return curr[RandomInt(0, int64(rand.Intn(len(curr))))]
}
