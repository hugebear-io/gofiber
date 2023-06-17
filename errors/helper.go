package errors

import (
	"github.com/bytedance/sonic"
)

func recast(from, to interface{}) error {
	switch v := from.(type) {
	case []byte:
		return sonic.Unmarshal(v, to)
	default:
		buf, err := sonic.Marshal(from)
		if err != nil {
			return err
		}

		return sonic.Unmarshal(buf, to)
	}
}

func ParsePostgresError(err error) error {
	obj := map[string]interface{}{}
	if err := recast(err, &obj); err != nil {
		return err
	}

	if code, ok := obj["Code"].(string); ok {
		if val, ok := errorMap[code]; ok {
			return val
		}
	}
	return err
}

var errorMap = map[string]error{
	"23505": NewBadRequestError("unique violation"),
	"23502": NewBadRequestError("not null violation"),
	"23503": NewBadRequestError("foreign key violation"),
	"42703": NewBadRequestError("undefined column"),
	"42P01": NewBadRequestError("undefined table"),
	"23514": NewBadRequestError("check violation"),
	"22001": NewBadRequestError("string data right truncation"),
	"42701": NewBadRequestError("duplicate column"),
	"22003": NewBadRequestError("numeric value out of range"),
	"22007": NewBadRequestError("invalid datetime format"),
	"42P09": NewBadRequestError("ambiguous alias"),
	"42P18": NewBadRequestError("indeterminate datatype"),
	"22023": NewBadRequestError("invalid parameter value"),
	"42P02": NewBadRequestError("undefined parameter"),
}
