package memory

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/caelifer/runner/service/store"
)

// Exported errors.
var (
	// ErrNotFound error is returned when object is not found in the data store.
	ErrNotFound = errors.New("object not found")
)

// memoryStore is an internal type that implements store.Service interface.
type memoryStore struct {
	entropy io.Reader
	logger  *log.Logger
}

// logRec is an internal type used for structured logging.
type logRec struct {
	Service   string `json:"service"`
	Operation string `json:"operation"`
	ID        string `json:"id,omitempty"`
	Success   bool   `json:"success,omitempty"`
	Error     string `json:"error,omitempty"`
	Duration  string `json:"duration"`
}

// String serializes structured log entry to JSON encoded string.
func (l logRec) String() string {
	out, _ := json.Marshal(&l)
	return string(out)
}

// New creates new MySQL based data store service.
func New() store.Service {
	logger := log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)

	return &memoryStore{
		logger:  logger,
		entropy: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Create new record in data store.
func (ms *memoryStore) Create(record store.Record) (err error) {
	defer func(t0 time.Time) {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		ms.logger.Printf("%v",
			logRec{
				Service:   "memory",
				Operation: "create",
				ID:        record.ID(),
				Error:     errStr,
				Duration:  fmt.Sprintf("%v", time.Since(t0)),
			},
		)
	}(time.Now())

	// run create op

	return
}

// Update existing record in data store.
func (ms *memoryStore) Update(id string, record store.Record) (err error) {
	defer func(t0 time.Time) {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		ms.logger.Printf("%v",
			logRec{
				Service:   "memory",
				Operation: "update",
				ID:        record.ID(),
				Success:   record.Success(),
				Error:     errStr,
				Duration:  fmt.Sprintf("%v", time.Since(t0)),
			},
		)
	}(time.Now())

	// Check if object exists first
	if !ms.isPresent(id) {
		err = ErrNotFound
		return
	}

	// Update state
	return
}

// Delete existing record from data store.
func (ms *memoryStore) Delete(id string) (err error) {
	defer func(t0 time.Time) {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		ms.logger.Printf("%v",
			logRec{
				Service:   "memory",
				Operation: "delete",
				ID:        id,
				Error:     errStr,
				Duration:  fmt.Sprintf("%v", time.Since(t0)),
			},
		)
	}(time.Now())

	// Check if object exists first
	if !ms.isPresent(id) {
		err = ErrNotFound
		return
	}

	// Update state
	return
}

// Get retrieves record from data store based on provided id.
func (ms *memoryStore) Get(id string) (record store.Record, err error) {
	defer func(t0 time.Time) {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		ms.logger.Printf("%v",
			logRec{
				Service:   "memory",
				Operation: "get",
				ID:        id,
				Error:     errStr,
				Duration:  fmt.Sprintf("%v", time.Since(t0)),
			},
		)
	}(time.Now())

	// get by id

	return
}

// GetAll fetches all records from data store as a slice.
func (ms *memoryStore) GetAll() (records []store.Record, err error) {
	defer func(t0 time.Time) {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		ms.logger.Printf("%+v",
			logRec{
				Service:   "memory",
				Operation: "get-all",
				Error:     errStr,
				Duration:  fmt.Sprintf("%v", time.Since(t0)),
			},
		)
	}(time.Now())

	return
}

// isPresent checks if record with given id exists in data store.
func (ms *memoryStore) isPresent(_ string) bool {
	return true
}
