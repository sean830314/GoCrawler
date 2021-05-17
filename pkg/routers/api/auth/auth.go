package auth

import (
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/sean830314/GoCrawler/pkg/app"
	"github.com/sean830314/GoCrawler/pkg/auth"
	"github.com/sean830314/GoCrawler/pkg/consts"
	"github.com/sean830314/GoCrawler/pkg/httputil"
	e "github.com/sean830314/GoCrawler/pkg/httputil"
	model "github.com/sean830314/GoCrawler/pkg/model/admin"
	"github.com/sean830314/GoCrawler/pkg/repository"
	admin "github.com/sean830314/GoCrawler/pkg/service/admin"
	"github.com/sean830314/GoCrawler/pkg/utils"
	"github.com/sirupsen/logrus"
)

var tokenManager = auth.TokenManager{}

type LoginForm struct {
	UserAccount  string `form:"userAccount" valid:"Required;MaxSize(100)"`
	UserPassword string `form:"userPassword" valid:"Required;MaxSize(100)"`
}

// @Summary Get User
// @Tags Auth
// @Produce  json
// @Accept multipart/form-data
// @Param userAccount formData string true "user account"
// @Param userPassword formData string true "user password"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/auth/login [post]
func Login(c *gin.Context) {
	var redisAuthService = auth.NewRedisAuthService()
	appG := app.Gin{C: c}
	var form LoginForm
	httpCode, errCode := app.BindAndValid(c, &form)
	if errCode != httputil.SUCCESS {
		appG.Response(httpCode, errCode, nil)
		return
	}
	dbConfig := repository.Config{
		Host:        utils.Env("GO_CRAWLER_DB_HOST", consts.DefDBHost),
		Port:        utils.Env("GO_CRAWLER_DB_PORT", consts.DefDBPort),
		User:        utils.Env("GO_CRAWLER_DB_USER", consts.DefDBUser),
		Pass:        utils.Env("GO_CRAWLER_DB_PASSWORD", consts.DefDBPass),
		Name:        utils.Env("GO_CRAWLER_DB_DBNAME", consts.DefDBName),
		SSLMode:     utils.Env("GO_CRAWLER_DB_SSLMODEL", consts.DefDBSSLMode),
		SSLCert:     utils.Env("GO_CRAWLER_DB_SSLCERT", consts.DefDBSSLCert),
		SSLKey:      utils.Env("GO_CRAWLER_DB_SSLKEY", consts.DefDBSSLKey),
		SSLRootCert: utils.Env("GO_CRAWLER_DB_SSLROOTCERT", consts.DefDBSSLRootCert),
	}
	db, err := repository.Connect(dbConfig)
	if err != nil {
		logrus.Error("Mysql connect error: ", err)
		appG.Response(http.StatusBadGateway, e.ERROR, err.Error())
	}
	repo := repository.New(db)
	svc := admin.NewBasicUserService(repo.User)
	req := model.UserReq{
		UserAccount:  &form.UserAccount,
		UserPassword: &form.UserPassword,
	}
	res, err := svc.Get(c, &req)
	if err != nil {
		logrus.Error("Get user error: ", err)
		appG.Response(http.StatusUnauthorized, e.INVALID_PARAMS, err.Error())
		return
	}
	ts, err := tokenManager.CreateToken(res.ID, res.UserAccount)
	if err != nil {
		logrus.Error("Create token error: ", err)
		appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, err.Error())
		return
	}
	err = redisAuthService.CreateAuth(res.ID, ts)
	if err != nil {
		logrus.Error("Store token to redis error: ", err)
		appG.Response(http.StatusUnprocessableEntity, e.UNAUTHORIZED, err.Error())
	}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	appG.Response(http.StatusOK, e.SUCCESS, tokens)
}

// @Summary Logout
// @Tags Auth
// @security Bearer
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/auth/logout [get]
func Logout(c *gin.Context) {
	var redisAuthService = auth.NewRedisAuthService()
	appG := app.Gin{C: c}
	au, err := auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		logrus.Error("Store token to redis error: ", err)
		appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, err.Error())
		return
	}
	err = redisAuthService.DeleteAuth(au)
	if err != nil { //if any goes wrong
		logrus.Error("Delete token to redis error: ", err)
		appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, err.Error())
		return
	}
	appG.Response(http.StatusOK, e.SUCCESS, "Successfully logged out")
}

type RefreshToken struct {
	RefreshToken string `json:"refreshToken" valid:"Required;MaxSize(100)"`
}

// @Summary Refresh
// @Tags Auth
// @security Bearer
// @Produce  json
// @Param refreshToken body RefreshToken true "refresh token"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/auth/refresh [post]
func Refresh(c *gin.Context) {
	appG := app.Gin{C: c}
	var redisAuthService = auth.NewRedisAuthService()
	mapToken := RefreshToken{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken.RefreshToken
	logrus.Info(refreshToken)
	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		logrus.Error("Parse error: ", err)
		appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, err.Error())
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		logrus.Error("unauthorized error: ", err)
		appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, err.Error())
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		userId, roleOk := claims["user_id"].(string)
		if roleOk == false {
			logrus.Error("unauthorized error: ", err)
			appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, err.Error())
			return
		}
		au, err := auth.ExtractTokenMetadata(c.Request)
		if err != nil {
			logrus.Error("Store token to redis error: ", err)
			appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, err.Error())
			return
		}
		err = redisAuthService.DeleteAuth(au)
		if err != nil { //if any goes wrong
			logrus.Error("Delete token to redis error: ", err)
			appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, err.Error())
			return
		}
		//Create new pairs of refresh and access tokens

		dbConfig := repository.Config{
			Host:        utils.Env("GO_CRAWLER_DB_HOST", consts.DefDBHost),
			Port:        utils.Env("GO_CRAWLER_DB_PORT", consts.DefDBPort),
			User:        utils.Env("GO_CRAWLER_DB_USER", consts.DefDBUser),
			Pass:        utils.Env("GO_CRAWLER_DB_PASSWORD", consts.DefDBPass),
			Name:        utils.Env("GO_CRAWLER_DB_DBNAME", consts.DefDBName),
			SSLMode:     utils.Env("GO_CRAWLER_DB_SSLMODEL", consts.DefDBSSLMode),
			SSLCert:     utils.Env("GO_CRAWLER_DB_SSLCERT", consts.DefDBSSLCert),
			SSLKey:      utils.Env("GO_CRAWLER_DB_SSLKEY", consts.DefDBSSLKey),
			SSLRootCert: utils.Env("GO_CRAWLER_DB_SSLROOTCERT", consts.DefDBSSLRootCert),
		}
		db, err := repository.Connect(dbConfig)
		if err != nil {
			logrus.Error("Mysql connect error: ", err)
			appG.Response(http.StatusBadGateway, e.ERROR, err.Error())
		}
		repo := repository.New(db)
		svc := admin.NewBasicUserService(repo.User)
		res, err := svc.GetById(c, userId)
		if err != nil {
			logrus.Error("Get user error: ", err)
			appG.Response(http.StatusUnauthorized, e.INVALID_PARAMS, err.Error())
			return
		}
		ts, createErr := tokenManager.CreateToken(res.ID, res.UserAccount)
		if createErr != nil {
			logrus.Error("Create token error: ", err)
			appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, createErr.Error())
			return
		}
		err = redisAuthService.CreateAuth(res.ID, ts)
		if err != nil {
			logrus.Error("Store token to redis error: ", err)
			appG.Response(http.StatusUnprocessableEntity, e.UNAUTHORIZED, err.Error())
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		appG.Response(http.StatusOK, e.SUCCESS, tokens)
	} else {
		appG.Response(http.StatusUnauthorized, e.UNAUTHORIZED, "refresh expired")
	}
}
