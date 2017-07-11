package store

type Record interface {
	ID() string
	Success() bool
}

type Service interface {
	Create(rec Record) error
	Update(id string, rec Record) error
	Delete(id string) error
	Get(id string) (rec Record, err error)
	GetAll() (recs []Record, err error)
}
