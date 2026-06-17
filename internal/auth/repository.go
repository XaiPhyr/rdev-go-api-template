package auth

import (
	"context"

	"github.com/XaiPhyr/rdev-go-api-template/internal/shared/models"
	"github.com/uptrace/bun"
)

type repository struct {
	db *bun.DB
}

func NewAuthRepository(db *bun.DB) *repository {
	return &repository{db: db}
}

func (r *repository) GetUsernameOrEmail(ctx context.Context, username string) (*models.User, error) {
	var user = new(models.User)

	err := r.db.NewSelect().
		Model(user).
		Where("username = ?", username).
		WhereOr("email = ?", username).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *repository) Register(ctx context.Context, user *models.User) error {
	_, err := r.db.NewInsert().Model(user).Exec(ctx)

	return err
}

func (r *repository) CheckUserPermission(ctx context.Context, userID int64, roleName string) ([]string, error) {
	var allPerms []string

	query := `
		WITH RECURSIVE role_hierarchy AS (
			SELECT role_id FROM user_roles WHERE user_id = ?
			UNION
			SELECT gr.role_id FROM user_groups ug
			JOIN group_roles gr ON ug.group_id = gr.group_id
			WHERE ug.user_id = ?
			UNION
			SELECT r.parent_id
			FROM roles r
			INNER JOIN role_hierarchy rh ON r.id = rh.role_id
			WHERE r.parent_id IS NOT NULL
		)
		SELECT r.name 
		FROM role_hierarchy rh
		JOIN roles r ON r.id = rh.role_id
		WHERE r.name = 'super_admin'
		UNION ALL
		SELECT DISTINCT p.slug
		FROM role_hierarchy rh
		JOIN role_permissions rp ON rp.role_id = rh.role_id
		JOIN permissions p ON p.id = rp.permission_id
	`

	err := r.db.NewRaw(query, userID, userID).Scan(ctx, &allPerms)

	if err != nil {
		return nil, err
	}

	return allPerms, nil
}
