// endpoints/user.go
package endpoints

import (
	"github.com/gin-gonic/gin"
	"github.com/norrsign/rullafy-data-api/repo"
	"github.com/vanern/goapi/framework/crud"
)

// RegisterUserEndpoints mounts /users and wires in pagination.
func RegisterUserEndpoints(rg *gin.RouterGroup) {
	r := repo.GetUserRepoSingleton()
	crud.RegisterCRUD(
		rg,
		r.List,
		r.Get,
		r.Create,
		r.Update,
		r.Delete,
	)
}
