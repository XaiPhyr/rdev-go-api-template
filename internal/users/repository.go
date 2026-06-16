package users

import (
	"context"
	"time"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/dto"
	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
	"github.com/uptrace/bun"
)

type repository struct {
	db *bun.DB
}

func NewUserRepository(db *bun.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().
		Model(&user).
		Exec(ctx)

	return err
}

func (r *repository) ReadOne(ctx context.Context, uuid string) (*models.User, error) {
	user := new(models.User)

	err := r.db.NewSelect().
		Model(user).
		Where("uuid = ?", uuid).
		Scan(ctx)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) ReadAll(ctx context.Context, q dto.Query) ([]models.User, int, error) {
	var users []models.User
	search := "%" + q.Search + "%"

	total, err := r.db.NewSelect().
		Model(&users).
		Limit(q.Limit).
		Offset(q.Offset).
		Order(q.OrderBy).
		WhereGroup("OR", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.WhereOr("u.first_name ILIKE ?", search).
				WhereOr("u.last_name ILIKE ?", search).
				WhereOr("u.email ILIKE ?", search)
		}).
		ScanAndCount(ctx)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (r *repository) Update(ctx context.Context, user *models.User) error {
	cols := []string{"first_name", "last_name", "addr1", "addr2", "city", "postal"}

	_, err := r.db.NewUpdate().
		Model(&user).
		Column(cols...).
		Set("u.updated_at = ?", time.Now()).
		WherePK().
		Exec(ctx)

	return err
}

func (r *repository) Delete(ctx context.Context, uuid string) error {
	_, err := r.db.NewDelete().
		Model((*models.User)(nil)).
		Where("uuid = ?", uuid).
		Exec(ctx)

	return err
}
