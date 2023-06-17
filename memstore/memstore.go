package memstore

import "time"

type Memstore interface {
	Get(key string, val interface{}) error
	Set(key string, val interface{}, duration *time.Duration) error
	Delete(key string) error
}
