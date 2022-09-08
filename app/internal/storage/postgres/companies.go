package postgres

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
)

type CompaniesQueries interface {
	CreateCompany(ctx context.Context, params Company) (*uint64, error)
	DeleteCompany(ctx context.Context, compID uint64) error
	GetCompanies(ctx context.Context, params Company) ([]*Company, error)
	GetCompanyByID(ctx context.Context, compID uint64) (*Company, error)
	UpdateCompany(ctx context.Context, compID uint64, data Company) (*uint64, error)
}

type Company struct {
	ID      uint64
	Name    string
	Code    string
	Country string
	Website string
	Phone   string
}

func (q *Queries) CreateCompany(
	ctx context.Context,
	data Company,
) (*uint64, error) {
	builder := q.builder.
		Insert("companies").
		SetMap(
			map[string]any{
				`name`:    data.Name,
				`code`:    data.Code,
				`country`: data.Country,
				`website`: data.Website,
				`phone`:   data.Phone,
			},
		).
		Suffix(` RETURNING id`)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	var id uint64
	err = q.tx.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return &id, nil
}

func (q *Queries) GetCompanies(
	ctx context.Context,
	params Company,
) ([]*Company, error) {
	res := []*Company{}
	builder := q.builder.
		Select(
			`id`,
			`name`,
			`code`,
			`country`,
			`website`,
			`phone`,
		).
		From("companies")
	builder = makeGetWheres(builder, params)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	err = pgxscan.Select(ctx, q.tx, &res, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return res, nil
}

func makeGetWheres(builder sq.SelectBuilder, params Company) sq.SelectBuilder {
	if params.Code != "" {
		builder = builder.Where(sq.Like{`code`: fmt.Sprintf("%%%s%%", params.Code)})
	}
	if params.Name != "" {
		builder = builder.Where(sq.Like{`name`: fmt.Sprintf("%%%s%%", params.Name)})
	}
	if params.Website != "" {
		builder = builder.Where(sq.Like{`website`: fmt.Sprintf("%%%s%%", params.Website)})
	}
	if params.Country != "" {
		builder = builder.Where(sq.Like{`country`: fmt.Sprintf("%%%s%%", params.Country)})
	}
	if params.Phone != "" {
		builder = builder.Where(sq.Like{`phone`: fmt.Sprintf("%%%s%%", params.Phone)})
	}
	return builder
}

func (q *Queries) UpdateCompany(
	ctx context.Context,
	compID uint64,
	data Company,
) (*uint64, error) {
	builder := q.builder.
		Update("companies").
		SetMap(
			map[string]any{
				`name`:    data.Name,
				`code`:    data.Code,
				`country`: data.Country,
				`website`: data.Website,
				`phone`:   data.Phone,
			},
		).
		Where(sq.Eq{`id`: compID}).
		Suffix(` RETURNING id`)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}
	var id uint64
	err = q.tx.QueryRow(ctx, query, args...).Scan(&id)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return &id, nil
}

func (q *Queries) GetCompanyByID(
	ctx context.Context,
	compID uint64,
) (*Company, error) {
	builder := q.builder.
		Select(
			`id`,
			`name`,
			`code`,
			`country`,
			`website`,
			`phone`,
		).
		From("companies").
		Where(sq.Eq{`id`: compID})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	var res Company
	err = pgxscan.Get(ctx, q.tx, &res, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query: %w", err)
	}

	return &res, nil
}

func (q *Queries) DeleteCompany(
	ctx context.Context,
	compID uint64,
) error {
	builder := q.builder.
		Delete("companies").
		Where(sq.Eq{`id`: compID})

	query, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("build query: %w", err)
	}

	_, err = q.tx.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return nil
}
