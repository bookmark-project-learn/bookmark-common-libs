package response

import (
	"errors"
	"net/http"

	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

var (
	BadRequestError    = errors.New("Invalid request")
	InternalError      = errors.New("Internal error")
	ConflictError      = errors.New("Conflict error")
	NotFoundError      = errors.New("Not found")
	MapErrorToHttpCode = map[error]int{
		BadRequestError: http.StatusBadRequest,
		InternalError:   http.StatusInternalServerError,
		ConflictError:   http.StatusConflict,
		NotFoundError:   http.StatusNotFound,
	}
)

func ErrorHandling(err error) error {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return NotFoundError
	case errors.Is(err, gorm.ErrDuplicatedKey), errors.Is(err, gorm.ErrForeignKeyViolated):
		return ConflictError
	case errors.Is(err, gorm.ErrInvalidTransaction) ||
		errors.Is(err, gorm.ErrMissingWhereClause) ||
		errors.Is(err, gorm.ErrInvalidValue) ||
		errors.Is(err, gorm.ErrInvalidData):
		return BadRequestError
	default:
		log.Error().Err(err)
		return InternalError
	}
}
