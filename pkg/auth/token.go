package auth

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid"
)

type AccessDetails struct {
	TokenUuid   string
	UserId      string
	UserAccount string
}
type TokenDetails struct {
	AccessToken  string
	RefreshToken string
	Authorized   bool
	TokenUuid    string
	RefreshUuid  string
	AtExpires    int64
	RtExpires    int64
}

type TokenManager struct{}

func NewTokenService() *TokenManager {
	return &TokenManager{}
}

type TokenInterface interface {
	CreateToken(userId, userName string) (*TokenDetails, error)
}

//Token implements the TokenInterface
var _ TokenInterface = &TokenManager{}

func (t *TokenManager) CreateToken(userId, userAccount string) (*TokenDetails, error) {
	td := &TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 30).Unix() //expires after 30 min
	td.TokenUuid = uuid.NewV4().String()
	td.Authorized = true

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.TokenUuid + "++" + userId

	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{
		"authorized":   td.Authorized,
		"access_uuid":  td.TokenUuid,
		"user_id":      userId,
		"user_account": userAccount,
		"exp":          td.AtExpires,
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	//Creating Refresh Token
	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = td.TokenUuid + "++" + userId

	rtClaims := jwt.MapClaims{
		"authorized":   td.Authorized,
		"access_uuid":  td.RefreshUuid,
		"user_id":      userId,
		"user_account": userAccount,
		"exp":          td.RtExpires,
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)

	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	return td, nil
}

func TokenValid(r *http.Request) error {
	token, err := VerifyToken(r)
	if err != nil {
		return err
	}
	_, ok := token.Claims.(jwt.Claims)
	if !ok && !token.Valid {
		return err
	}
	return nil
}

//verify token function
func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	if tokenString == "" {
		return nil, errors.New("Get access token failed.")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

//get the token from the request header
func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func Extract(token *jwt.Token) (*AccessDetails, error) {

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		userId, userOk := claims["user_id"].(string)
		userAccount, userAccountOk := claims["user_account"].(string)
		if ok == false || userOk == false || userAccountOk == false {
			return nil, errors.New("unauthorized")
		} else {
			return &AccessDetails{
				TokenUuid:   accessUuid,
				UserId:      userId,
				UserAccount: userAccount,
			}, nil
		}
	}
	return nil, errors.New("something went wrong")
}

func ExtractTokenMetadata(r *http.Request) (*AccessDetails, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return nil, err
	}
	acc, err := Extract(token)
	if err != nil {
		return nil, err
	}
	return acc, nil
}
