package service

import (
	"context"
	"time"

	gonanoid "github.com/matoous/go-nanoid"
	"github.com/opentracing/opentracing-go"
	model "github.com/sean830314/GoCrawler/pkg/model/admin"
)

// RoleService describes the service.
type RoleService interface {
	// [method=get,expose=true,router=items]
	List(ctx context.Context) (res []*model.RoleRes, err error)
	// [method=post,expose=true,router=items]
	Add(ctx context.Context, role *model.RoleReq) (res *model.RoleRes, err error)
	// [method=put,expose=true,router=items/:id]
	Update(ctx context.Context, id string, role *model.RoleReq) (res *model.RoleRes, err error)
	// [method=delete,expose=true,router=items/:id]
	Delete(ctx context.Context, id string) (err error)
}

type basicRoleService struct {
	repo model.RoleRepository
}

func (b *basicRoleService) List(ctx context.Context) (res []*model.RoleRes, err error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "List")
	defer span.Finish()
	// TODO ＋ the business logic of List
	res = make([]*model.RoleRes, 0)

	rr, err := b.repo.List(ctx)
	if err != nil {
		return
	}
	for _, r := range rr {
		item := model.RoleRes(*r)
		res = append(res, &item)
	}
	return
}
func (b *basicRoleService) Add(ctx context.Context, role *model.RoleReq) (res *model.RoleRes, err error) {
	// TODO implement the business logic of Add
	id, _ := gonanoid.ID(21)

	t := new(model.Role)
	t.ID = id
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	if role.Name != nil {
		t.Name = *role.Name
	}
	if role.Slug != nil {
		t.Slug = *role.Slug
	}
	if err := b.repo.Add(ctx, t); err != nil {
		return res, err
	}
	x := model.RoleRes(*t)
	return &x, nil
}
func (b *basicRoleService) Update(ctx context.Context, id string, role *model.RoleReq) (res *model.RoleRes, err error) {
	// TODO implement the business logic of Update
	dt, err := b.repo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	dt.UpdatedAt = time.Now()
	if role.Name != nil {
		dt.Name = *role.Name
	}
	if role.Slug != nil {
		dt.Slug = *role.Slug
	}

	if err := b.repo.Update(ctx, dt); err != nil {
		return nil, err
	}
	x := model.RoleRes(*dt)
	return &x, nil
}
func (b *basicRoleService) Delete(ctx context.Context, id string) (err error) {
	// TODO implement the business logic of Delete
	return b.repo.Delete(ctx, id)
	// return err
}

// NewBasicRoleService returns a naive, stateless implementation of RoleService.
func NewBasicRoleService(repo model.RoleRepository) RoleService {
	return &basicRoleService{repo: repo}
}
