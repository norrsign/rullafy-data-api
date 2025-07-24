package repo

import (
	"context"
	"strconv"

	"github.com/norrsign/rullafy-data-api/db"
	"github.com/norrsign/rullafy-data-api/db/models"
	"github.com/vanern/goapi/types"
)

type UserRepoStruct struct {
	Q *models.Queries
}

var userRepo *UserRepoStruct = nil

func GetUserRepoSingleton() *UserRepoStruct {
	if userRepo != nil {
		return userRepo
	}
	if db.Qrs == nil {
		panic("dn.InitDB is not called and the Queries is nil")
	}
	userRepo = &UserRepoStruct{Q: db.Qrs}
	return userRepo
}

// List returns a paginated list of users plus total count.
func (r *UserRepoStruct) List(ctx context.Context, limit, offset int32) (types.ListResult[models.User], error) {

	// calculate offset

	// fetch the page
	users, err := r.Q.ListUsers(ctx, models.ListUsersParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return types.ListResult[models.User]{}, err
	}

	// fetch total count
	total, err := r.Q.CountUsers(ctx)
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

func (r *UserRepoStruct) Get(ctx context.Context, id int64) (models.User, error) {
	return r.Q.GetUser(ctx, strconv.FormatInt(id, 10))
}

func (r *UserRepoStruct) Create(ctx context.Context, params models.CreateUserParams) (models.User, error) {
	return r.Q.CreateUser(ctx, params)
}

func (r *UserRepoStruct) Update(ctx context.Context, id int64, params models.UpdateUserParams) (models.User, error) {
	params.ID = strconv.FormatInt(id, 10)
	return r.Q.UpdateUser(ctx, params)
}

func (r *UserRepoStruct) Delete(ctx context.Context, id int64) error {
	return r.Q.DeleteUser(ctx, strconv.FormatInt(id, 10))
}
