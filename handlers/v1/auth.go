package v1

import (
	"database/sql"
	"fit-byte/repositories"
	"fit-byte/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type AuthRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=32"`
}

type AuthHandler struct {
	Repo *repositories.UserRepository
}

type AuthLoginResponse struct {
	IdentityNumber string `json:"identityNumber"`
	Name           string `json:"name"`
}

type AuthRegisterResponse struct {
	IdentityNumber string `json:"identityNumber"`
	Name           string `json:"name"`
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{
		Repo: repositories.NewUserRepository(db),
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Repo.FindUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Email not found"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Password mismatch"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	err = h.Repo.UpdateTokenById(user.ID, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"email": user.Email,
		"token": token,
	})
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req AuthRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.Repo.FindUserByEmail(req.Email)
	if err != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hashing password"})
			return
		}

		user, err := h.Repo.CreateUser(req.Email, string(hashedPassword))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		token, err := utils.GenerateJWT(user.ID, user.Email)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		err = h.Repo.UpdateTokenById(user.ID, token)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"email": user.Email,
			"token": token,
		})
	} else {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}
}
