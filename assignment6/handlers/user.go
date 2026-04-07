package handlers

import (
	"assignment6/models"
	"assignment6/utils"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// Register godoc
// @Summary Register a new user
// @Description Create a new user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RegisterRequest true "Register Request"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 409 {object} map[string]interface{}
// @Router /register [post]
func Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.Mu.Lock()
	defer models.Mu.Unlock()

	if _, exists := models.Users[req.Email]; exists {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Simple account creation
	user := &models.User{
		ID:               models.NextID,
		Email:            req.Email,
		Password:         req.Password, // No hashing for simplicity, as it's an assignment
		Role:             models.RoleUser,
		IsVerified:       false,
		VerificationCode: "1234", // Fixed 4-digit code for testing the easy task
	}
	models.Users[req.Email] = user
	models.NextID++

	fmt.Printf("Verification code for %s: %s\n", user.Email, user.VerificationCode)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered. Please verify your email with code 1234.",
		"user":    user,
	})
}

// Login godoc
// @Summary Log in a user
// @Description Authenticate user and return tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body LoginRequest true "Login Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.Mu.RLock()
	user, exists := models.Users[req.Email]
	models.Mu.RUnlock()

	if !exists || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	accessToken, refreshToken, err := utils.GenerateTokens(user.ID, user.Email, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

// GetMe godoc
// @Summary Get current user info
// @Description Returns the email of the authenticated user
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/me [get]
func GetMe(c *gin.Context) {
	email, _ := c.Get("email")
	
	// Task 1: Return user email from JWT
	c.JSON(http.StatusOK, gin.H{
		"email": email,
	})
}

// PromoteUser godoc
// @Summary Promote a user to admin
// @Description Changes the role of a specific user to 'admin'
// @Tags users
// @Security ApiKeyAuth
// @Param id path int true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /users/promote/{id} [patch]
func PromoteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	models.Mu.Lock()
	defer models.Mu.Unlock()

	found := false
	for _, user := range models.Users {
		if user.ID == id {
			user.Role = models.RoleAdmin
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User promoted to admin successfully"})
}

type VerifyRequest struct {
	Email string `json:"email" binding:"required"`
	Code  string `json:"code" binding:"required"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// VerifyEmail godoc
// @Summary Verify email address
// @Description Validates the 4-digit code sent to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body VerifyRequest true "Verify Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /verify [post]
func VerifyEmail(c *gin.Context) {
	var req VerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// ...

	models.Mu.Lock()
	defer models.Mu.Unlock()

	user, exists := models.Users[req.Email]
	if !exists || user.VerificationCode != req.Code {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid code or user"})
		return
	}

	user.IsVerified = true
	c.JSON(http.StatusOK, gin.H{"message": "Account verified successfully"})
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Issues a new Access Token using a Refresh Token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body RefreshRequest true "Refresh Request"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /users/refresh [post]
func RefreshToken(c *gin.Context) {
	var req RefreshRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := utils.ValidateToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token: " + err.Error()})
		return
	}

	// Issue new Access Token
	accessToken, newRefreshToken, err := utils.GenerateTokens(claims.UserID, claims.Email, claims.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not regenerate tokens"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": newRefreshToken,
	})
}
