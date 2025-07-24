package repo

import (
	"context"
	"strconv"

	"github.com/norrsign/rullafy-data-api/db/models"
	"github.com/vanern/goapi/types"
)

type UserRepo struct {
	q *models.Queries
}

func NewUserRepo(q *models.Queries) *UserRepo {
	return &UserRepo{q: q}
}

// List returns a paginated list of users plus total count.
func (r *UserRepo) List(ctx context.Context, limit, offset int32) (types.ListResult[models.User], error) {

	// calculate offset

	// fetch the page
	users, err := r.q.ListUsers(ctx, models.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return types.ListResult[models.User]{}, err
	}

	// fetch total count
	total, err := r.q.CountUsers(ctx)
	if err != nil {
		return types.ListResult[models.User]{}, err
	}
	p := offset / limit
	if offset%limit != 0 {
		p++
	}
	return types.ListResult[models.User]{
		Data:  users,
		Page:  p,
		Total: total,
	}, nil
}

// Get, Create, Update, Delete remain exactly as before:

func (r *UserRepo) Get(ctx context.Context, id int64) (models.User, error) {
	return r.q.GetUser(ctx, strconv.FormatInt(id, 10))
}

func (r *UserRepo) Create(ctx context.Context, params models.CreateUserParams) (models.User, error) {
	return r.q.CreateUser(ctx, params)
}

func (r *UserRepo) Update(ctx context.Context, id int64, params models.UpdateUserParams) (models.User, error) {
	params.ID = strconv.FormatInt(id, 10)
	return r.q.UpdateUser(ctx, params)
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	return r.q.DeleteUser(ctx, strconv.FormatInt(id, 10))
}
