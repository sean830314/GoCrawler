package admin

import (
	"context"
	"time"
)

type User struct {
	ID           string    `gorm:"primaryKey" json:"id"`
	UserAccount  string    `json:"userAccount"`
	UserPassword string    `json:"userPassword"`
	Name         string    `json:"name"`
	NickName     string    `json:"nickName"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type UserRepository interface {
	Add(context.Context, *User) error
	Delete(context.Context, string) error
	Update(context.Context, *User) error
	List(context.Context) (res []*User, err error)
	Get(context.Context, *User) (res *User, err error)
	GetById(context.Context, string) (res *User, err error)
}

type UserReq struct {
	Name         *string `form:"name" valid:"Required;MaxSize(100)"`
	UserAccount  *string `form:"userAccount" valid:"Required;MaxSize(100)"`
	UserPassword *string `form:"userPassword" valid:"Required;MaxSize(100)"`
	NickName     *string `form:"nickName" valid:"Required;MaxSize(100)"`
	Role         *string `form:"role" valid:"Required;MaxSize(100)"`
}

type UserRes User
