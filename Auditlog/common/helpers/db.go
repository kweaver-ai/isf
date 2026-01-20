package helpers

import (
	"database/sql"
	"errors"
	"fmt"

	"AuditLog/gocommon/api"
)

//nolint:gocritic
func TxRollbackOrCommit(tx *sql.Tx, err *error, logger api.Logger) {
	re := recover()
	if *err != nil || re != nil {
		if re != nil {
			logger.Errorf("db panic: %v", re)

			// 如果是panic并且err为nil，将panic转换为error
			if *err == nil {
				//nolint:goerr113
				*err = fmt.Errorf("%v", re)
			}
		}

		if e := tx.Rollback(); e != nil {
			*err = e
		}
	} else {
		if e := tx.Commit(); e != nil {
			*err = e
		}
	}
}

func IsSqlNotFound(err error) bool {
	if err == nil {
		return false
	}

	return errors.Is(err, sql.ErrNoRows)
}

func CloseRows(rows *sql.Rows, logger api.Logger) {
	if rows != nil {
		if rowsErr := rows.Err(); rowsErr != nil {
			logger.Errorln(rowsErr)
		}

		if closeErr := rows.Close(); closeErr != nil {
			logger.Errorln(closeErr)
		}
	}
}
