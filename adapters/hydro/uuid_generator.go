package hydro

import (
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/boreq/errors"
	"github.com/oklog/ulid"
)

type UUIDGenerator struct {
	entropy io.Reader
	mutex   sync.Mutex
}

func NewUUIDGenerator() *UUIDGenerator {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)

	return &UUIDGenerator{
		entropy: entropy,
	}
}

func (u *UUIDGenerator) Generate() (string, error) {
	u.mutex.Lock()
	defer u.mutex.Unlock()

	uuid, err := ulid.New(ulid.Timestamp(time.Now()), u.entropy)
	if err != nil {
		return "", errors.Wrap(err, "could not create a new ulid")
	}

	return uuid.String(), nil
}
