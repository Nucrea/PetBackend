package repos

import (
	"backend/internal/integrations"
	"context"
	"database/sql"
	"errors"
	"time"

	"go.opentelemetry.io/otel/trace"
)

type ShortlinkDTO struct {
	Id        string
	Url       string
	ExpiresAt time.Time
}

type ShortlinkRepo interface {
	AddShortlink(ctx context.Context, dto ShortlinkDTO) error
	GetShortlink(ctx context.Context, id string) (*ShortlinkDTO, error)
	DeleteExpiredShortlinks(ctx context.Context, limit int) (int, error)
}

func NewShortlinkRepo(db integrations.SqlDB, tracer trace.Tracer) ShortlinkRepo {
	return &shortlinkRepo{db, tracer}
}

type shortlinkRepo struct {
	db     integrations.SqlDB
	tracer trace.Tracer
}

func (u *shortlinkRepo) AddShortlink(ctx context.Context, dto ShortlinkDTO) error {
	_, span := u.tracer.Start(ctx, "postgres::AddShortlink")
	defer span.End()

	query := `insert into shortlinks (id, url, expires_at) values ($1, $2, $3);`
	_, err := u.db.ExecContext(ctx, query, dto.Id, dto.Url, dto.ExpiresAt)
	return err
}

func (u *shortlinkRepo) GetShortlink(ctx context.Context, id string) (*ShortlinkDTO, error) {
	_, span := u.tracer.Start(ctx, "postgres::GetShortlink")
	defer span.End()

	query := `select url, expires_at from shortlinks where id = $1;`
	row := u.db.QueryRowContext(ctx, query, id)
	if err := row.Err(); err != nil {
		return nil, err
	}

	dto := &ShortlinkDTO{Id: id}
	err := row.Scan(&dto.Url, &dto.ExpiresAt)
	if err == nil {
		return dto, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}

	return nil, err
}

func (u *shortlinkRepo) DeleteExpiredShortlinks(ctx context.Context, limit int) (int, error) {
	_, span := u.tracer.Start(ctx, "postgres::CheckExpiredShortlinks")
	defer span.End()

	query := `
	select count(*) from (
		delete from shortlinks
		where id in (
			select id
				from shortlinks
				where current_date > expires_at
				limit $1
		)
		returning *
	);`
	row := u.db.QueryRowContext(ctx, query, limit)
	if err := row.Err(); err != nil {
		return 0, err
	}

	count := 0
	err := row.Scan(&count)
	if err == nil {
		return count, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}

	return 0, err
}
