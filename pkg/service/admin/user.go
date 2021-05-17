package service

import (
	"context"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/opentracing/opentracing-go"
	model "github.com/sean830314/GoCrawler/pkg/model/admin"
)

// UserService describes the service.
type UserService interface {
	// [method=get,expose=true,router=items]
	GetById(ctx context.Context, id string) (res *model.UserRes, err error)
	// [method=get,expose=true,router=items]
	Get(ctx context.Context, user *model.UserReq) (res *model.UserRes, err error)
	// [method=get,expose=true,router=items]
	List(ctx context.Context) (res []*model.UserRes, err error)
	// [method=post,expose=true,router=items]
	Add(ctx context.Context, user *model.UserReq) (res *model.UserRes, err error)
	// [method=put,expose=true,router=items/:id]
	Update(ctx context.Context, id string, user *model.UserReq) (res *model.UserRes, err error)
	// [method=delete,expose=true,router=items/:id]
	Delete(ctx context.Context, id string) (err error)
}

type basicUserService struct {
	repo model.UserRepository
}

func (b *basicUserService) GetById(ctx context.Context, id string) (res *model.UserRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get")
	defer span.Finish()
	r, err := b.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}
	userRes := model.UserRes(*r)
	return &userRes, nil
}

func (b *basicUserService) Get(ctx context.Context, user *model.UserReq) (res *model.UserRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "Get")
	defer span.Finish()
	u := new(model.User)
	if user.UserAccount != nil {
		u.UserAccount = *user.UserAccount
	}
	if user.UserPassword != nil {
		u.UserPassword = *user.UserPassword
	}
	r, err := b.repo.Get(ctx, u)
	if err != nil {
		return nil, err
	}
	userRes := model.UserRes(*r)
	return &userRes, nil
}

func (b *basicUserService) List(ctx context.Context) (res []*model.UserRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "List")
	defer span.Finish()
	// TODO ＋ the business logic of List
	res = make([]*model.UserRes, 0)

	rr, err := b.repo.List(ctx)
	if err != nil {
		return
	}
	for _, r := range rr {
		item := model.UserRes(*r)
		res = append(res, &item)
	}
	return
}

func (b *basicUserService) Add(ctx context.Context, user *model.UserReq) (res *model.UserRes, err error) {
	// TODO implement the business logic of Add
	id, _ := gonanoid.ID(21)

	t := new(model.User)
	t.ID = id
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	if user.UserAccount != nil {
		t.UserAccount = *user.UserAccount
	}
	if user.UserPassword != nil {
		t.UserPassword = *user.UserPassword
	}
	if user.Name != nil {
		t.Name = *user.Name
	}
	if user.NickName != nil {
		t.NickName = *user.NickName
	}
	if user.Role != nil {
		t.Role = *user.Role
	}
	if err := b.repo.Add(ctx, t); err != nil {
		return res, err
	}
	x := model.UserRes(*t)
	return &x, nil
}

func (b *basicUserService) Update(ctx context.Context, id string, user *model.UserReq) (res *model.UserRes, err error) {
	// TODO implement the business logic of Update
	dt, err := b.repo.GetById(ctx, id)
	if err != nil {
		return nil, err
	}

	dt.UpdatedAt = time.Now()
	if user.Name != nil {
		dt.Name = *user.Name
	}
	if user.NickName != nil {
		dt.NickName = *user.NickName
	}
	if user.UserAccount != nil {
		dt.UserAccount = *user.UserAccount
	}
	if user.UserPassword != nil {
		dt.UserPassword = *user.UserPassword
	}
	if user.Role != nil {
		dt.Role = *user.Role
	}
	if err := b.repo.Update(ctx, dt); err != nil {
		return nil, err
	}
	x := model.UserRes(*dt)
	return &x, nil
}

func (b *basicUserService) Delete(ctx context.Context, id string) (err error) {
	// TODO implement the business logic of Delete
	return b.repo.Delete(ctx, id)
	// return err
}

// NewBasicUserService returns a naive, stateless implementation of UserService.
func NewBasicUserService(repo model.UserRepository) UserService {
	return &basicUserService{repo: repo}
}
