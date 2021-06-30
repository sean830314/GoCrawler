package admin

import (
	"context"
	"time"
)

type Role struct {
	ID        string    `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"unique" json:"name"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type RoleRepository interface {
	Add(context.Context, *Role) error
	Delete(context.Context, string) error
	Update(context.Context, *Role) error
	List(context.Context) (res []*Role, err error)
	Get(context.Context, string) (res *Role, err error)
}

type RoleReq struct {
	Name *string `form:"name" valid:"Required;MaxSize(100)"`
	Slug *string `form:"slug" valid:"Required;MaxSize(100)"`
}

type RoleRes Role
