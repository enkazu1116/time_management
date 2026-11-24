package repository

import (
	"uuid"

	"github.com/jackc/pgx/v5/pgtype"
)

/*
 * Go 標準の UUID を Postgres 向け pgtype.UUID に変換する
 */
func pgtypeConvertUUID(id uuid.UUID) pgtype.UUID {
	return pgtype.UUID{
		Bytes: id,
		Valid: true,
	}
}
