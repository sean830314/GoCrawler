package repository

import (
	"context"
	"sync"

	model "github.com/sean830314/GoCrawler/pkg/model/admin"
	"gorm.io/gorm"
)

var _ model.UserRepository = (*userRepository)(nil)

type userRepository struct {
	mu sync.RWMutex
	db *gorm.DB
}

func (repo *userRepository) Get(ctx context.Context, userID string) (res *model.User, err error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	res = new(model.User)
	err = repo.db.WithContext(ctx).Where("id", userID).Find(res).Error
	return
}

func (repo *userRepository) Add(ctx context.Context, user *model.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if err := repo.db.WithContext(ctx).Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) Delete(ctx context.Context, userID string) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	if err := repo.db.WithContext(ctx).Delete(&model.User{ID: userID}).Error; err != nil {
		return err
	}
	return nil
}

func (repo *userRepository) Update(ctx context.Context, user *model.User) error {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	result := repo.db.WithContext(ctx).Model(&model.User{ID: user.ID}).UpdateColumns(
		map[string]interface{}{
			"name":       user.Name,
			"nick_name":  user.NickName,
			"role":       user.Role,
			"updated_at": user.UpdatedAt,
		},
	)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (repo *userRepository) List(ctx context.Context) (res []*model.User, err error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()

	err = repo.db.WithContext(ctx).Order("created_at desc").Find(&res).Error
	return
}

func NewUserRepository(db *gorm.DB) model.UserRepository {
	return &userRepository{
		mu: sync.RWMutex{},
		db: db,
	}
}
