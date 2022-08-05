package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	uuid "github.com/nu7hatch/gouuid"
)

func GenerateShortCode() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999

	return strconv.Itoa(rand.Intn((max - min) + min))
}

func GenerateMediumCode() string {
	rand.Seed(time.Now().UnixNano())
	min := 1000000000
	max := 9999999999

	return strconv.Itoa(rand.Intn((max - min) + min))
}

func GenerateTransactionCode(transctionType string) string {
	rand.Seed(time.Now().UnixNano())
	min := 100000
	max := 999999

	t := time.Now()
	time := t.Format("20060102150405")
	randomCode := strconv.Itoa(rand.Intn((max - min) + min))
	result := fmt.Sprintf("trx:%v:%v:%v", time, randomCode, transctionType)
	return result
}

func GenerateUUID() string {
	u, _ := uuid.NewV4()
	return u.String()
}
