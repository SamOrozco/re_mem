package hash

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/segmentio/fasthash/fnv1a"
)

func NewHash(value string) uint64 {
	return fnv1a.HashString64(value)
}

func NewHashString(value string) string {
	return fmt.Sprint(NewHash(value))
}

func NewRandomKey() string {
	return uuid.New().String()
}
