package vo

import (
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/provider/errors"
	"github.com/google/uuid"
)

type useUuid struct {
	value string
}

func (t useUuid) String() string {
	return t.value
}

func GenerateUuid[T ~struct{ useUuid }]() T {
	return T{useUuid{uuid.NewString()}}
}

// ConstructUuid Hack constructor biar "method" disini bisa di "extend" sama value object Id lainnya
func newUuid[T ~struct{ useUuid }](str string, customError ...errors.InvariantError) (T, error) {
	_, err := uuid.Parse(str)
	if err != nil {
		if len(customError) > 0 {
			return T{}, fmt.Errorf("%w but got %s instead", customError[0], str)
		} else {
			return T{}, fmt.Errorf("%w but got %s instead", err, str)
		}
	}
	return T{useUuid{str}}, nil
}
