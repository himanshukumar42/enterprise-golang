package controllers

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/himanshukumar42/enterprise/forms"
	"github.com/himanshukumar42/enterprise/models"
)

type AuthController struct{}

var authModel = new(models.AuthModel)

func (ctrl AuthController) TokenValid(c *gin.Context) {
	tokenAuth, err := authModel.ExtractTokenMetaData(c.Request)
	if err != nil {
		// Token either expired or not valid
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "please login first"})
		return
	}
	userID, err := authModel.FetchAuth(tokenAuth)
	if err != nil {
		// Token does not exists in Redis (User logged out or expired)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "please login first"})
		return
	}

	// To be called from GetUserID()
	c.Set("userID", userID)
}

func (ctrl AuthController) Refresh(c *gin.Context) {
	var tokenForm forms.Token

	if c.ShouldBindJSON(&tokenForm) != nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"message": "invalid form", "form": tokenForm})
		c.Abort()
		return
	}

	// verify the token
	token, err := jwt.Parse(tokenForm.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		// Make sure that the token method confirms to "SigningMethodHMAC"
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	fmt.Println("Token: ", token)
	// if there is an error the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization, please login again"})
		return
	}

	fmt.Println("2 error: ", err)
	// is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization, please login again"})
		return
	}
	// Since token is valid, get the uuid
	claims, ok := token.Claims.(jwt.MapClaims) // the token claims should conform to MapClaims
	fmt.Println("Claims: ", claims)
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) // convert the interface to string
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization, please login again"})
			return
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization, please login again"})
			return
		}

		// Delete the previous refresh token
		fmt.Println("refreshUUID: ", refreshUUID)
		deleted, delErr := authModel.DeleteAuth(refreshUUID)
		fmt.Println("delErr: , deleted ", delErr, deleted)
		if delErr != nil || deleted == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization, please login again"})
			return
		}
		fmt.Println("idhara to aaya")

		// create new pairs of refresh and access tokens
		ts, createErr := authModel.CreateToken(userID)
		if createErr != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "invalid authorization, pleas login again"})
			return
		}

		// save the tokens metadata to redis
		saveErr := authModel.CreateAuth(userID, ts)
		if saveErr != nil {
			c.JSON(http.StatusForbidden, gin.H{"message": "invalid authorization, pleaes login again"})
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusOK, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid authorization, please login again"})
	}
}
