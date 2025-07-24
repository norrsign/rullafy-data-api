// endpoints/user.go
package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/norrsign/rullafy-data-api/db/models"
	"github.com/norrsign/rullafy-data-api/repo"
	"github.com/vanern/goapi/framework/crud"
)

// RegisterUserEndpoints mounts /users and wires in pagination.
func RegisterUserEndpoints(rg *gin.RouterGroup, ur *repo.UserRepo) {
	crud.RegisterCRUD[models.User, models.CreateUserParams, models.UpdateUserParams](
		rg.Group("/users"),
		ur.List,
		ur.Get,
		ur.Create, 
		ur.Update,
		ur.Delete,
	)
}
