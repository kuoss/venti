package handler

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/kuoss/common/logger"
	"github.com/kuoss/venti/pkg/handler/api"
	"github.com/kuoss/venti/pkg/model"
	userService "github.com/kuoss/venti/pkg/service/user"
	"gorm.io/gorm"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
	// todo service to database
	userService *userService.UserService
}

func NewAuthHandler(s *userService.UserService) *authHandler {
	return &authHandler{s}
}

func (h *authHandler) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" {
		api.ResponseError(c, api.ErrorUnauthorized, fmt.Errorf("username is empty"))
		return
	}
	if password == "" {
		api.ResponseError(c, api.ErrorUnauthorized, fmt.Errorf("password is empty"))
		return
	}

	user, err := h.userService.FindByUsername(username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			api.ResponseError(c, api.ErrorUnauthorized, fmt.Errorf("username not found"))
			return
		}
		api.ResponseError(c, api.ErrorUnauthorized, fmt.Errorf("FindByUsername err: %w", err))
		return
	}

	if !checkPassword(password, user.Hash) {
		logger.Infof("User login failed.")
		api.ResponseError(c, api.ErrorUnauthorized, fmt.Errorf("username or password is incorrect"))
		return
	}

	user = issueToken(user)
	err = h.userService.Save(user)
	if err != nil {
		logger.Errorf("update token err: %s", err.Error())
		api.ResponseError(c, api.ErrorInternal, fmt.Errorf("token save err: %w", err))
		return
	}

	logger.Infof("user '%s' logged in successfully.", user.Username)
	c.JSON(200, gin.H{
		"message":  "You are logged in.",
		"token":    user.Token,
		"userID":   user.ID,
		"username": user.Username,
	})
}

func checkPassword(plain string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

func issueToken(user model.User) model.User {
	if user.Token != "" && user.TokenExpires.After(time.Now()) {
		user.TokenExpires = time.Now().Add(48 * time.Hour)
	} else {
		b := make([]byte, 16)
		_, _ = rand.Read(b)
		token := fmt.Sprintf("%x", b)
		user.Token = token
		user.TokenExpires = time.Now().Add(48 * time.Hour)
	}
	return user
}

func (h *authHandler) Logout(c *gin.Context) {
	//deleteTokenIfWeCan(c)
	// token delete if we can
	tokenFromHeader := c.GetHeader("Authorization")
	if tokenFromHeader == "" || !strings.HasPrefix(tokenFromHeader, "Bearer ") {
		// todo normal logout?
		return
	}
	tokenFromHeader = strings.TrimPrefix(tokenFromHeader, "Bearer ")
	userID := c.GetHeader("UserID")

	user, err := h.userService.FindByUserIdAndToken(userID, tokenFromHeader)
	if err != nil {
		return
	}

	user.Token = ""
	user.TokenExpires = time.Now().Add(-480 * time.Hour)
	err = h.userService.Save(user)
	if err != nil {
		return
	}
	c.JSON(200, gin.H{
		"message": "You are logged out.",
	})
}

// TODO: add to handler routes check if we have to this on method
// only checks userid, Bearer tokene exist. not validation.
func (h *authHandler) HeaderRequired(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "valid token required",
		})
		return
	}

	userID := c.GetHeader("UserID")
	if userID == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "userId required",
		})
		return
	}
	c.Next()
}
