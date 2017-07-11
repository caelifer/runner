package mysqlstore

import (
	"errors"
	"io"
	"math/rand"
	"sync"
	"time"

	"github.com/caelifer/runner/service/store"
)

var (
	ErrNotFound = errors.New("object not found")
)

// type that implements store.Service service interface
type mysqlstoredummy struct {
	sync.RWMutex
	store   map[string]store.Record
	entropy io.Reader
}

func New() store.Service {
	return &mysqlstoredummy{
		store:   make(map[string]store.Record),
		entropy: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (ms *mysqlstoredummy) Create(record store.Record) (err error) {
	t0 := time.Now()
	defer func(t0 time.Time) {
		logger.Log(
			"service", "mysqlstoredummy",
			"operation", "create",
			"id", record.ID(),
			"error", err,
			"duration", time.Since(t0),
		)
	}(t0)

	// Update state
	ms.Lock()
	ms.store[record.ID()] = record
	ms.Unlock()

	return
}

func (ms *mysqlstoredummy) Update(id string, record store.Record) (err error) {
	defer func(t0 time.Time) {
		logger.Log(
			"service", "mysqlstoredummy",
			"operation", "update",
			"id", id,
			"success", record.Success(),
			"error", err,
			"duration", time.Since(t0),
		)
	}(time.Now())

	// Check if object exists first
	if !ms.isPresent(id) {
		err = ErrNotFound
		return
	}

	// Update state
	ms.Lock()
	ms.store[id] = record
	ms.Unlock()

	return
}

func (ms *mysqlstoredummy) Delete(id string) (err error) {
	defer func(t0 time.Time) {
		logger.Log(
			"service", "mysqlstoredummy",
			"operation", "delete",
			"id", id,
			"error", err,
			"duration", time.Since(t0),
		)
	}(time.Now())

	// Check if object exists first
	if !ms.isPresent(id) {
		err = ErrNotFound
		return
	}

	// Update state
	ms.Lock()
	delete(ms.store, id)
	ms.Unlock()

	return
}

func (ms *mysqlstoredummy) Get(id string) (record store.Record, err error) {
	defer func(t0 time.Time) {
		logger.Log(
			"service", "mysqlstoredummy",
			"operation", "get",
			"id", id,
			"error", err,
			"duration", time.Since(t0),
		)
	}(time.Now())

	ms.RLock()
	defer ms.RUnlock()

	var ok bool
	if record, ok = ms.store[id]; !ok {
		err = ErrNotFound
	}

	return
}

func (ms *mysqlstoredummy) GetAll() (records []store.Record, err error) {
	defer func(t0 time.Time) {
		logger.Log(
			"service", "mysqlstoredummy",
			"operation", "getAll",
			"error", err,
			"duration", time.Since(t0),
		)
	}(time.Now())

	ms.RLock()
	defer ms.RUnlock()

	// Collect all job objects
	for _, v := range ms.store {
		records = append(records, v)
	}

	return
}

func (ms *mysqlstoredummy) isPresent(id string) bool {
	ms.RLock()
	defer ms.RUnlock()
	_, ok := ms.store[id]
	return ok
}
