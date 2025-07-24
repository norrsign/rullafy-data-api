package crud

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/vanern/goapi/types"
)

// TODO: set the limit to a sensible default, e.g. 100
const maxInt32 = "2147483647"

// RegisterCRUD wires up List/Get/Create/Update/Delete handlers
// on the provided group. T is your model type,
// C is the create‑params type, U is the update‑params type.
//
// The list func is now expected to honor pagination:
//
//	func(ctx context.Context, limit, offset int32) ([]T, error)
func RegisterCRUD[T any, C any, U any](
	group *gin.RouterGroup,
	list func(ctx context.Context, limit, offset int32) (types.ListResult[T], error),
	get func(ctx context.Context, id int64) (T, error),
	create func(ctx context.Context, params C) (T, error),
	update func(ctx context.Context, id int64, params U) (T, error),
	del func(ctx context.Context, id int64) error,
) {
	// ─── LIST with pagination ───────────────────────────────────────────
	group.GET("", func(c *gin.Context) {
		// parse ?limit= and ?offset=, with defaults
		limitStr := c.DefaultQuery("limit", maxInt32)
		offsetStr := c.DefaultQuery("offset", "0")

		lim, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || lim < 1 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be between 1 and 100"})
			return
		}
		off, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}

		items, err := list(c.Request.Context(), int32(lim), int32(off))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, items)
	})

	// ─── GET by ID ───────────────────────────────────────────────────────
	group.GET("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		item, err := get(c.Request.Context(), id)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	// ─── CREATE ─────────────────────────────────────────────────────────
	group.POST("", func(c *gin.Context) {
		var input C
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item, err := create(c.Request.Context(), input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, item)
	})

	// ─── UPDATE ─────────────────────────────────────────────────────────
	group.PUT("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		var input U
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		item, err := update(c.Request.Context(), id, input)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, item)
	})

	// ─── DELETE ─────────────────────────────────────────────────────────
	group.DELETE("/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}
		if err := del(c.Request.Context(), id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.Status(http.StatusNoContent)
	})
}
