package helper

import (
	"errors"
	"net/http"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// unfinished
func PGErrorToHTTPCode(err error) int {
	var res int

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.CaseNotFound:
			res = http.StatusNotFound
		case pgerrcode.UniqueViolation:
			res = http.StatusConflict
		default:
			res = http.StatusInternalServerError
		}
	}

	return res
}
