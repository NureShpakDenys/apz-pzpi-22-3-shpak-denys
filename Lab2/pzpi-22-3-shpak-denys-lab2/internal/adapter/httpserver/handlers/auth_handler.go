package handlers // import "wayra/internal/adapter/httpserver/handlers"

import (
	"context"
	"net/http"
	"time"
	"wayra/internal/core/domain/dtos"
	"wayra/internal/core/domain/models"
	"wayra/internal/core/port/services"

	"github.com/gin-gonic/gin"
)

// AuthHandler is a handler for auth endpoints
type AuthHandler struct {
	authService services.AuthService // Service for auth operations
	tokenExpiry time.Duration        // Token expiry time
	LogService  services.LogService  // Service for logging actions
	UserService services.UserService // Service for user operations
}

// NewAuthHandler creates a new AuthHandler
// authService: Service for auth operations
// tokenExpiry: Token expiry time
// Returns: A new AuthHandler
func NewAuthHandler(
	authService services.AuthService,
	tokenExpiry time.Duration,
	logService services.LogService,
	userService services.UserService,
) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		tokenExpiry: tokenExpiry,
		LogService:  logService,
		UserService: userService,
	}
}

// AuthCredentials is the request for the RegisterUser and LoginUser endpoints
type AuthCredentials struct {
	// Username is the name of the user
	// Example: john_doe
	Username string `json:"username" example:"john_doe"`

	// Password is the password of the user
	// Example: password123
	Password string `json:"password" example:"password123"`
}

// RegisterUser godoc
// @Summary      Register a new user
// @Description  Registers a new user with the provided details
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user body AuthCredentials true "User details"
// @Router       /auth/register [post]
func (h *AuthHandler) RegisterUser(c *gin.Context) {
	var request AuthCredentials
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.LogService.LogAction(nil, "register", "Failed register attempt: "+err.Error(), false)
		return
	}

	if request.Username == "" || len(request.Username) < 4 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username must be at least 4 characters long"})
		h.LogService.LogAction(nil, "register", "Failed register attempt: Username must be at least 4 characters long", false)
		return
	}

	if request.Password == "" || len(request.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password must be at least 8 characters long"})
		h.LogService.LogAction(nil, "register", "Failed register attempt: Password must be at least 8 characters long", false)
		return
	}

	user := models.User{
		Name:     request.Username,
		Password: request.Password,
		RoleID:   UserRole,
	}

	if err := h.authService.RegisterUser(context.Background(), &user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to register user"})
		h.LogService.LogAction(nil, "register", "Failed register attempt: "+err.Error(), false)
		return
	}
	h.LogService.LogAction(&user.ID, "register", "User registered successfully", true)
	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

// LoginUser godoc
// @Summary      Login user
// @Description  Authenticates a user and returns a token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body AuthCredentials true "User credentials"
// @Router       /auth/login [post]
func (h *AuthHandler) LoginUser(c *gin.Context) {
	var credentials AuthCredentials

	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		h.LogService.LogAction(nil, "login", "Failed login: "+err.Error(), false)
		return
	}

	token, err := h.authService.LoginUser(context.Background(), credentials.Username, credentials.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		h.LogService.LogAction(nil, "login", "Failed login: "+err.Error(), false)
		return
	}

	users, _ := h.UserService.Where(context.Background(), &models.User{Name: credentials.Username})
	user := users[0]

	userToReturn, err := h.UserService.GetByID(context.Background(), user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user details"})
		h.LogService.LogAction(nil, "login", "Failed to get user details: "+err.Error(), false)
		return
	}

	type UserWithRole struct {
		dtos.UserDTO
		Role string `json:"role"`
	}

	userDTO := UserWithRole{
		UserDTO: dtos.UserDTO{
			ID:   userToReturn.ID,
			Name: userToReturn.Name,
		},
		Role: userToReturn.Role.Name,
	}

	c.SetCookie("token", token, int(h.tokenExpiry.Seconds()), "/", "localhost", false, true)
	h.LogService.LogAction(&user.ID, "login", "User logged in successfully", true)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  userDTO,
	})
}

// LogoutUser godoc
// @Summary      Logout user
// @Description  Logs out a user by invalidating their token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Router       /auth/logout [post]
func (h *AuthHandler) LogoutUser(c *gin.Context) {
	userID, err := GetUserIDFromToken(c)
	if err != nil {
		userID = nil
	}

	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})

	h.LogService.LogAction(userID, "logout", "User logged out successfully", true)
}
