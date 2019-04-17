package generator

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

var (
	entropy = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func NewID() string {
	return ulid.MustNew(ulid.Timestamp(time.Now()), entropy).String()
}
