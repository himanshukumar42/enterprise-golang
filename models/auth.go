package models

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/himanshukumar42/enterprise/db"
)

type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	AccessUUID   string
	RefreshUUID  string
	AtExpires    int64
	RtExpires    int64
}

type AccessDetails struct {
	AccessUUID string
	UserID     int64
}

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type AuthModel struct{}

func (m AuthModel) CreateToken(userID int64) (*TokenDetails, error) {
	td := &TokenDetails{
		AccessUUID:  uuid.New().String(),
		RefreshUUID: uuid.New().String(),
		AtExpires:   time.Now().Add(time.Hour * 2).Unix(),
		RtExpires:   time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	var err error

	// Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUUID
	atClaims["user_id"] = userID
	atClaims["exp"] = td.AtExpires

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	// Creating Refresh Token
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUUID
	rtClaims["user_id"] = userID
	rtClaims["expires"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return td, err
}

func (m AuthModel) CreateAuth(userID int64, td *TokenDetails) error {
	at := time.Unix(td.AtExpires, 0) // converting Unix to UTC (to time object)
	rt := time.Unix(td.RtExpires, 0)
	now := time.Now()

	errAccess := db.GetRedis().Set(td.AccessUUID, strconv.Itoa(int(userID)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := db.GetRedis().Set(td.RefreshUUID, strconv.Itoa(int(userID)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}
	return nil
}

// Extract Token from Authorization HTTP Request
func (m AuthModel) ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	// normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

// Verify Extracted Token
func (m AuthModel) VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := m.ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Make sure that the token method confirms to "SigningMethodHMAC"
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

// Token Valid
func (m AuthModel) TokenValid(r *http.Request) error {
	token, err := m.VerifyToken(r)
	if err != nil {
		return err
	}
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		return err
	}
	return nil
}

// ExtractTokenMetaData
func (m AuthModel) ExtractTokenMetaData(r *http.Request) (*AccessDetails, error) {
	token, err := m.VerifyToken(r)
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUUID, ok := claims["access_uuid"].(string)
		if !ok {
			return nil, err
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			return nil, err
		}
		return &AccessDetails{
			AccessUUID: accessUUID,
			UserID:     userID,
		}, nil
	}
	return nil, err
}

func (m AuthModel) FetchAuth(authD *AccessDetails) (int64, error) {
	userid, err := db.GetRedis().Get(authD.AccessUUID).Result()
	if err != nil {
		return 0, err
	}
	userID, _ := strconv.ParseInt(userid, 10, 64)
	return userID, nil
}

func (m AuthModel) DeleteAuth(givenUUID string) (int64, error) {
	deleted, err := db.GetRedis().Del(givenUUID).Result()
	if err != nil {
		return 0, err
	}
	return deleted, nil
}
