package vo

import (
	"github.com/Marcellinom/tenant-management-saas/provider/errors"
	"github.com/google/uuid"
)

type UseUuid struct {
	value string
}

func (t UseUuid) String() string {
	return t.value
}

func GenerateUuid[T ~struct{ UseUuid }]() T {
	return T{UseUuid{uuid.NewString()}}
}

// ConstructUuid Hack constructor biar "method" disini bisa di "extend" sama value object Id lainnya
func newUuid[T ~struct{ UseUuid }](str string, customError ...errors.InvariantError) (T, error) {
	_, err := uuid.Parse(str)
	if err != nil {
		if len(customError) > 0 {
			return T{}, customError[0]
		} else {
			return T{}, err
		}
	}
	return T{UseUuid{str}}, nil
}
