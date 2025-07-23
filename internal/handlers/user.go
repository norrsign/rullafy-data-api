package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/norrsign/rullafy-data-api/internal/middleware"
	"github.com/norrsign/rullafy-data-api/internal/repo"
	"github.com/norrsign/rullafy-data-api/internal/typs"
)

func RegisterUser(r *gin.Engine, ur *repo.UserRepo) {
	grp := r.Group("/users")
	grp.Use(middleware.JWTAuth())

	grp.GET("", func(c *gin.Context) {
		users, err := ur.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	})

	grp.POST("", func(c *gin.Context) {
		var body struct{ Job string }
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		au, _ := typs.GetAuthUser(c)
		user, err := ur.Create(au.ID, body.Job)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	})

	grp.GET("/:id", func(c *gin.Context) {
		u, err := ur.Get(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, u)
	})

	grp.PUT("/:id", func(c *gin.Context) {
		var body struct{ Job string }
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		u, err := ur.Update(c.Param("id"), body.Job)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, u)
	})

	grp.DELETE("/:id", func(c *gin.Context) {
		if err := ur.Delete(c.Param("id")); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
