package handler

import (
	"crypto/rand"
	"fmt"
	"github.com/kuoss/venti/pkg/model"
	"github.com/kuoss/venti/pkg/store"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type authHandler struct {
	// todo service to database
	service *store.UserStore
}

func NewAuthHandler(us *store.UserStore) *authHandler {
	return &authHandler{us}
}

func (ah *authHandler) Login(c *gin.Context) {

	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(200, gin.H{
			"errors": map[string]string{
				"common": "The username or password is invalid.",
			},
		})
		return
	}

	user, err := ah.service.FindByUsername(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":    "error",
			"errorType": "internal error",
			"error":     "something wrong with db", // todo
		})
		return
	}

	if checkPassword(password, user.Hash) {
		user := issueToken(user)
		err := ah.service.Save(user)
		if err != nil {
			log.Printf("token info update to database failed %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"status":    "error",
				"errorType": "internal error",
				"error":     "something wrong with db", // todo error define
			})
			return
		}

		log.Printf("User '%s' logged in successfully.\n", user.Username)
		c.JSON(200, gin.H{
			"message":  "You are logged in.",
			"token":    user.Token,
			"userID":   user.ID,
			"username": user.Username,
		})
		log.Println("User login ok.")
		return
	}

	log.Println("User login failed.")

	// todo status code 422?
	c.JSON(422, gin.H{
		"status":    "error",
		"errorType": "authentication_failed",
		"error":     "The username or password is incorrect",
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

func (ah *authHandler) Logout(c *gin.Context) {
	//deleteTokenIfWeCan(c)
	// token delete if we can
	tokenFromHeader := c.GetHeader("Authorization")
	if tokenFromHeader == "" || !strings.HasPrefix(tokenFromHeader, "Bearer ") {
		// todo normal logout?
		return
	}
	tokenFromHeader = strings.TrimPrefix(tokenFromHeader, "Bearer ")
	userID := c.GetHeader("UserID")

	user, err := ah.service.FindByUserIdAndToken(userID, tokenFromHeader)
	if err != nil {
		return
	}

	user.Token = ""
	user.TokenExpires = time.Now().Add(-480 * time.Hour)
	err = ah.service.Save(user)
	if err != nil {
		return
	}
	c.JSON(200, gin.H{
		"message": "You are logged out.",
	})
}

// TODO: add to handler routes check if we have to this on moethod
// only checks userid, Bearer tokene exist. not validation.
func (ah *authHandler) HeaderRequired(c *gin.Context) {
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
