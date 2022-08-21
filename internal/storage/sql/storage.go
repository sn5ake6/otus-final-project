package sqlstorage

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type Storage struct {
	dsn string
	db  *sqlx.DB
	ctx context.Context
}

func New(dsn string) *Storage {
	return &Storage{
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Open("pgx", s.dsn)
	if err != nil {
		return err
	}

	s.db = db
	s.ctx = ctx

	return nil
}

func (s *Storage) Close() error {
	return s.db.Close()
}

func (s *Storage) AddToBlacklist(subnet string) error {
	sql := `
		INSERT INTO blacklist
		  (subnet)
		VALUES
		  ($1)
		;
	`

	_, err := s.db.ExecContext(s.ctx, sql, subnet)

	return err
}

func (s *Storage) DeleteFromBlacklist(subnet string) error {
	sql := `
		DELETE
		FROM
		  blacklist
		WHERE
		  subnet = $1
		;
	`
	_, err := s.db.ExecContext(s.ctx, sql, subnet)

	return err
}

func (s *Storage) FindIPInBlacklist(ip string) (bool, error) {
	sql := `
		SELECT
		  subnet
		FROM
		  blacklist
		WHERE
		  subnet >> $1
		;
	`

	res, err := s.db.ExecContext(s.ctx, sql, ip)
	if err != nil {
		return false, err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

func (s *Storage) AddToWhitelist(subnet string) error {
	sql := `
		INSERT INTO whitelist
		  (subnet)
		VALUES
		  ($1)
		;
	`

	_, err := s.db.ExecContext(
		s.ctx,
		sql,
		subnet,
	)

	return err
}

func (s *Storage) DeleteFromWhitelist(subnet string) error {
	sql := `
		DELETE
		FROM
		whitelist
		WHERE
		  subnet = $1
		;
	`
	_, err := s.db.ExecContext(s.ctx, sql, subnet)

	return err
}

func (s *Storage) FindIPInWhitelist(ip string) (bool, error) {
	sql := `
		SELECT
		  subnet
		FROM
		  whitelist
		WHERE
		  subnet >> $1
		;
	`

	res, err := s.db.ExecContext(s.ctx, sql, ip)
	if err != nil {
		return false, err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}
