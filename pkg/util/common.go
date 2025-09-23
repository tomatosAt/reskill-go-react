package util

import (
	"net"
	"os"
	"reflect"
	"strconv"
	"time"
	"unsafe"

	"github.com/google/uuid"
)

func ParseDuration(t string) time.Duration {
	d, _ := time.ParseDuration(t)
	return d
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}
	_, ok := set[item]
	return ok
}

func AtoI(s string, v int) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return v
	}
	return i
}

func AtoF(s string, v float64) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return v
	}
	return f
}

func IPv4Tester(ip string) bool {
	return net.ParseIP(ip) != nil
}

/*
Prefer the best performance Byte to String converter
*/

func B2S(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func S2B(s string) (b []byte) {
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return b
}

func GenUniqueIdV7() uuid.UUID {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		return uuid.New()
	}
	return uuidV7
}
