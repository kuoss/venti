package server

import (
	"crypto/rand"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func login(c *gin.Context) {
	var form LoginForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(422, gin.H{
			"errors": map[string]string{
				"common": "The username or password is invalid.",
			},
		})
		return
	}
	var user User
	err := db.First(&user, "username = ?", "admin").Error
	if err == nil {
		if checkPassword(form.Password, user.Hash) {
			token := issueToken(user)
			log.Printf("User '%s' logged in successfully.", user.Username)
			c.JSON(200, gin.H{
				"message":  "You are logged in.",
				"token":    token,
				"userID":   user.ID,
				"username": user.Username,
			})
			return
		}
	}
	log.Printf("User login failed.")
	c.JSON(422, gin.H{
		"status":    "error",
		"errorType": "authentication_failed",
		"error":     "The username or password is incorrect",
	})
}

func checkPassword(plain string, hashed string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain)) == nil
}

func logout(c *gin.Context) {
	deleteTokenIfWeCan(c)
	c.JSON(200, gin.H{
		"message": "You are logged out.",
	})
}

func deleteTokenIfWeCan(c *gin.Context) {
	token := c.GetHeader("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")
	userID := c.GetHeader("UserID")
	var user User
	err := db.First(&user, "ID = ? AND token = ?", userID, token).Error
	if err != nil {
		return
	}
	log.Println("User '" + user.Username + "' has logged out.")
	user.Token = ""
	user.TokenExpires = time.Now().Add(-480 * time.Hour)
	db.Save(&user)
}

// TODO: add to api routes
func TokenRequired(c *gin.Context) {
	// log.Println("URL=", c.Request.URL)
	// if c.Request.URL.String() == "/api/datasources/targets" {
	// 	c.Next()
	// 	return
	// }

	token := c.GetHeader("Authorization")
	if token == "" || !strings.HasPrefix(token, "Bearer ") {
		c.String(http.StatusUnauthorized, `{"message":"invalid token"}`)
		c.Abort()
		return
	}
	token = strings.TrimPrefix(token, "Bearer ")
	userID := c.GetHeader("UserID")
	var user User
	err := db.First(&user, "id = ? AND token = ?", userID, token).Error
	if err != nil {
		c.String(http.StatusUnauthorized, `{"message":"cannot find token"}`)
		c.Abort()
		return
	}
	if user.TokenExpires.Before(time.Now()) {
		c.String(http.StatusUnauthorized, `{"message":"expired token"}`)
		c.Abort()
		return
	}
	c.Next()
}

func issueToken(user User) string {
	if user.Token != "" && user.TokenExpires.After(time.Now()) {
		user.TokenExpires = time.Now().Add(48 * time.Hour)
	} else {
		b := make([]byte, 16)
		rand.Read(b)
		token := fmt.Sprintf("%x", b)
		user.Token = token
		user.TokenExpires = time.Now().Add(48 * time.Hour)
	}
	db.Save(&user)
	return user.Token
}
