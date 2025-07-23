package repo

import (
	"context"

	"github.com/norrsign/rullafy-data-api/db"
)

type UserRepo struct{ q *db.Queries }

func NewUserRepo(q *db.Queries) *UserRepo { return &UserRepo{q: q} }

func (r *UserRepo) List() ([]db.User, error) {
	return r.q.ListUsers(context.Background())
}
func (r *UserRepo) Get(id string) (db.User, error) {
	return r.q.GetUser(context.Background(), id)
}
func (r *UserRepo) Create(id, job string) (db.User, error) {
	return r.q.CreateUser(context.Background(), db.CreateUserParams{ID: id, Job: job})
}
func (r *UserRepo) Update(id, job string) (db.User, error) {
	return r.q.UpdateUser(context.Background(), db.UpdateUserParams{ID: id, Job: job})
}
func (r *UserRepo) Delete(id string) error {
	return r.q.DeleteUser(context.Background(), id)
}
