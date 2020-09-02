package resources

import (
	"sync"

	"github.com/paulbellamy/ratecounter"
	"go.uber.org/zap"
)

type DB struct {
	mx  sync.Mutex
	m   map[string]*ratecounter.RateCounter
	Log *zap.SugaredLogger
}

func NewDB(logger *zap.SugaredLogger) *DB {
	m := make(map[string]*ratecounter.RateCounter)
	return &DB{
		m:   m,
		Log: logger,
	}
}

func (db *DB) Close() {
	for k := range db.m {
		delete(db.m, k)
	}
}

func (db *DB) Load(key string) (*ratecounter.RateCounter, bool) {
	db.mx.Lock()
	defer db.mx.Unlock()
	val, ok := db.m[key]
	return val, ok
}

func (db *DB) GetAllKeys() []string {
	db.mx.Lock()
	defer db.mx.Unlock()
	keys := make([]string, 0, len(db.m))
	for k := range db.m {
		keys = append(keys, k)
	}
	return keys
}

func (db *DB) Store(key string, value *ratecounter.RateCounter) {
	db.mx.Lock()
	defer db.mx.Unlock()
	db.m[key] = value
}
