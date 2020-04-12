package meet

import (
	"fmt"
	"sync/atomic"
)

type UUIDGeneratorMock struct {
	i uint64
}

func NewUUIDGeneratorMock() *UUIDGeneratorMock {
	return &UUIDGeneratorMock{}
}

func (u *UUIDGeneratorMock) Generate() (string, error) {
	i := atomic.AddUint64(&u.i, 1)
	return fmt.Sprintf("uuid-%d", i), nil
}
