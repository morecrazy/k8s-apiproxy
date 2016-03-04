package object

import (
	"errors"
)

var (
	RecordNotFoundError      = errors.New("record not found")
	RecordAlreadyExistsError = errors.New("record already exists")
	RecordNoneAffectedError  = errors.New("no record affected")
	DbError                  = errors.New("database error")
)
