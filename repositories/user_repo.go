package repositories

import (
	"github.com/norrsign/rullafy-data-api/db"
	"github.com/vanern/goapi/framework/apperr"
	"github.com/vanern/goapi/framework/persistence"
	"github.com/vanern/goapi/typs"
)

type UserRepo struct{ q *db.Queries }

func NewUserRepo(q *db.Queries) *UserRepo { return &UserRepo{q: q} }

func (r *UserRepo) List(p typs.Page) (typs.ListResult[db.User], error) {
	users, err := r.q.ListUsers(persistence.Ctx())
	if err != nil {
		return typs.ListResult[db.User]{}, err
	}
	// naive paging in memory (demo); replace by SQL LIMIT/OFFSET.
	start, end := sliceBounds(len(users), p)
	return typs.ListResult[db.User]{
		Data:  users[start:end],
		Page:  p,
		Total: len(users),
	}, nil
}

func (r *UserRepo) Get(id string) (db.User, error) {
	u, err := r.q.GetUser(persistence.Ctx(), id)
	if err != nil {
		return db.User{}, apperr.New(apperr.ErrNotFound, "user not found", err)
	}
	return u, nil
}

func (r *UserRepo) Create(in db.CreateUserParams) (db.User, error) {
	return r.q.CreateUser(persistence.Ctx(), in)
}
func (r *UserRepo) Update(in db.UpdateUserParams) (db.User, error) {
	return r.q.UpdateUser(persistence.Ctx(), in)
}
func (r *UserRepo) Delete(id string) error {
	return r.q.DeleteUser(persistence.Ctx(), id)
}

// ---------- tiny util ------------------------------------------------------
func sliceBounds(n int, p typs.Page) (int, int) {
	start := (p.Page - 1) * p.Size
	if start > n {
		start = n
	}
	end := start + p.Size
	if end > n {
		end = n
	}
	return start, end
}
