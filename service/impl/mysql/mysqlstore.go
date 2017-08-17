package mysql

import (
	"errors"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"fmt"

	"encoding/json"

	"github.com/caelifer/runner/service/store"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	ErrNotFound = errors.New("object not found")
)

// type that implements store.Service service interface
type mysqlstoredummy struct {
	entropy io.Reader
	db      *gorm.DB
	logger  *log.Logger
}

type logrec struct {
	Service   string `json:"service"`
	Operation string `json:"operation"`
	ID        string `json:"id,omitempty"`
	Error     string `json:"error,omitempty"`
	Duration  string `json:"duration"`
}

func (l logrec) String() string {
	out, _ := json.Marshal(&l)
	return string(out)
}

func New() *mysqlstoredummy {
	logger := log.New(os.Stderr, "", log.Ldate|log.Lmicroseconds|log.Lshortfile)
	db, err := gorm.Open("mysql", "root@/test?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		logger.Fatal(err)
	}

	return &mysqlstoredummy{
		db:      db,
		logger:  logger,
		entropy: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

func (ms *mysqlstoredummy) Create(record store.Record) (err error) {
	t0 := time.Now()
	defer func(t0 time.Time) {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		ms.logger.Printf("%v",
			logrec{
				Service:   "mysqldummy",
				Operation: "create",
				ID:        record.ID(),
				Error:     errStr,
				Duration:  fmt.Sprintf("%v", time.Since(t0)),
			},
		)
	}(t0)

	// run create op

	return
}

func (ms *mysqlstoredummy) Update(id string, record store.Record) (err error) {
	defer func(t0 time.Time) {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		ms.logger.Printf("%v",
			logrec{
				Service:   "mysqldummy",
				Operation: "update",
				ID:        record.ID(),
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

func (ms *mysqlstoredummy) Delete(id string) (err error) {
	defer func(t0 time.Time) {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		ms.logger.Printf("%v",
			logrec{
				Service:   "mysqldummy",
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

func (ms *mysqlstoredummy) Get(id string) (record store.Record, err error) {
	defer func(t0 time.Time) {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		ms.logger.Printf("%v",
			logrec{
				Service:   "mysqldummy",
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

func (ms *mysqlstoredummy) GetAll() (records []store.Record, err error) {
	defer func(t0 time.Time) {
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		ms.logger.Printf("%+v",
			logrec{
				Service:   "mysqldummy",
				Operation: "get-all",
				Error:     errStr,
				Duration:  fmt.Sprintf("%v", time.Since(t0)),
			},
		)
	}(time.Now())

	return
}

func (ms *mysqlstoredummy) isPresent(id string) bool {
	return true
}
