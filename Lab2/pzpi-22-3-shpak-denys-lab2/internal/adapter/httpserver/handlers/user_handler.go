package handlers // import "wayra/internal/adapter/httpserver/handlers"

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	dtoMapper "github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserHandler is a struct to handle user requests
type UserHandler struct {
	userService services.UserService // is a service to handle user business logic
	db          *gorm.DB
}

// NewUserHandler is a function to create a new UserHandler
// userService: is a service to handle user business logic
// returns: a new UserHandler
func NewUserHandler(userService services.UserService, db *gorm.DB) *UserHandler {
	return &UserHandler{
		userService: userService,
		db:          db,
	}
}

// UpdateUserRequest is a struct to handle user update request
type UpdateUserRequest struct {
	// Name is the name of the user
	// Example: John Doe
	Name string `json:"name"`
}

// GetUser godoc
// @Summary Get a user
// @Description Get a user by ID
// @ID get-user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security     BearerAuth
// @Router /user/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetByID(context.Background(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userDTO := &dtos.UserDTO{}
	if err = dtoMapper.Map(userDTO, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userDTO)
}

// GetUsers godoc
// @Summary Get users
// @ID get-users
// @Tags user
// @Accept json
// @Produce json
// @Security     BearerAuth
// @Router /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.Where(context.Background(), &models.User{})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	type UserWithRole struct {
		dtos.UserDTO
		Role string `json:"role"`
	}

	var userDTOs []UserWithRole

	for _, user := range users {
		var role string
		if err := h.db.Raw("SELECT name FROM roles WHERE id = ?", user.RoleID).Scan(&role).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		userDTOs = append(userDTOs, UserWithRole{
			UserDTO: dtos.UserDTO{
				ID:       user.ID,
				Name:     user.Name,
				Password: user.Password,
			},
			Role: role,
		})
	}
	c.JSON(http.StatusOK, userDTOs)
}

// UpdateUser godoc
// @Summary Update a user
// @Description Update a user by ID
// @ID update-user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body UpdateUserRequest true "User data"
// @Security     BearerAuth
// @Router /user/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetByID(context.Background(), uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	currentUserID, err := getUserIDFromToken(c)
	fmt.Println(currentUserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	currentUser, err := h.userService.GetByID(context.Background(), *currentUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if currentUser.Role.Name != "admin" || currentUser.ID != user.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden"})
		return
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	user.Companies = nil
	user.Role = models.Role{}

	if err := h.userService.Update(context.Background(), user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userDto := &dtos.UserDTO{}
	if err = dtoMapper.Map(userDto, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, userDto)
}

// DeleteUser godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @ID delete-user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Security     BearerAuth
// @Router /user/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.Delete(context.Background(), uint(userID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
