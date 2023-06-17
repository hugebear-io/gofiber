package errors

import "errors"

func Equal(err, target error) bool {
	return errors.Is(err, target)
}
