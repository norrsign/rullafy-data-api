package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/you/rullafy-data-api/internal/middleware"
	"github.com/you/rullafy-data-api/internal/repo"
	"github.com/you/rullafy-data-api/internal/typs"
)

func RegisterUser(r *gin.Engine, ur *repo.UserRepo) {
	grp := r.Group("/users")
	// all user routes need a logged-in user
	grp.Use(middleware.JWTAuth())

	grp.GET("", listUsers(ur))
	grp.POST("", createUser(ur))
	grp.GET("/:id", getUser(ur))
	grp.PUT("/:id", updateUser(ur))
	grp.DELETE("/:id", deleteUser(ur))
}

func listUsers(ur *repo.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		users, err := ur.List()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func createUser(ur *repo.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		// JSON body: { "job": "..." }
		var body struct{ Job string }
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		// use JWT sub as ID if user exists
		au, _ := typs.GetAuthUser(c)
		id := au.ID
		user, err := ur.Create(id, body.Job)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func getUser(ur *repo.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		u, err := ur.Get(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}
		c.JSON(http.StatusOK, u)
	}
}

func updateUser(ur *repo.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
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
	}
}

func deleteUser(ur *repo.UserRepo) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := ur.Delete(c.Param("id")); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	}
}
