package repository

import (
	"context"
	"sync"

	model "github.com/sean830314/GoCrawler/pkg/model/admin"
	"gorm.io/gorm"
)

var _ model.RoleRepository = (*roleRepository)(nil)

type roleRepository struct {
	mu sync.RWMutex
	db *gorm.DB
}

func (repo *roleRepository) Get(ctx context.Context, roleID string) (res *model.Role, err error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	res = new(model.Role)
	err = repo.db.WithContext(ctx).Where("id", roleID).Find(res).Error
	return
}

func (repo *roleRepository) Add(ctx context.Context, role *model.Role) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if err := repo.db.WithContext(ctx).Create(role).Error; err != nil {
		return err
	}
	return nil
}

func (repo *roleRepository) Delete(ctx context.Context, roleID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if err := repo.db.WithContext(ctx).Delete(&model.Role{ID: roleID}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *roleRepository) Update(ctx context.Context, role *model.Role) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	result := repo.db.WithContext(ctx).Model(&model.Role{ID: role.ID}).UpdateColumns(
		map[string]interface{}{
			"name":       role.Name,
			"slug":       role.Slug,
			"updated_at": role.UpdatedAt,
		},
	)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *roleRepository) List(ctx context.Context) (res []*model.Role, err error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	err = repo.db.WithContext(ctx).Order("created_at desc").Find(&res).Error
	return
}

func NewRoleRepository(db *gorm.DB) model.RoleRepository {
	return &roleRepository{
		mu: sync.RWMutex{},
		db: db,
	}
}
