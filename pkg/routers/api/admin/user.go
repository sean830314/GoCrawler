package admin

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sean830314/GoCrawler/pkg/app"
	"github.com/sean830314/GoCrawler/pkg/consts"
	"github.com/sean830314/GoCrawler/pkg/httputil"
	e "github.com/sean830314/GoCrawler/pkg/httputil"
	model "github.com/sean830314/GoCrawler/pkg/model/admin"
	"github.com/sean830314/GoCrawler/pkg/repository"
	admin "github.com/sean830314/GoCrawler/pkg/service/admin"
	"github.com/sean830314/GoCrawler/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/unknwon/com"
)

// @Summary List Users
// @Produce  json
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/users [get]
func ListUsers(c *gin.Context) {
	appG := app.Gin{C: c}
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
		logrus.Error("err", err)
		os.Exit(1)
	}
	repo := repository.New(db)
	svc := admin.NewBasicUserService(repo.User)
	res, err := svc.List(c)
	if err != nil {
		logrus.Error("err", err)
	}
	appG.Response(http.StatusOK, e.SUCCESS, res)
}

type AddUserForm struct {
	Name     string `form:"name" valid:"Required;MaxSize(100)"`
	NickName string `form:"nickName" valid:"Required;MaxSize(255)"`
	Role     string `form:"role" valid:"Required;MaxSize(255)"`
}

// @Summary Add User
// @Produce  json
// @Accept multipart/form-data
// @Param name formData string true "name"
// @Param nickName formData string true "nick name"
// @Param role formData string true "role"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/users [post]
func AddUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var form AddUserForm
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
		logrus.Error("err", err)
		os.Exit(1)
	}
	repo := repository.New(db)
	svc := admin.NewBasicUserService(repo.User)
	req := model.UserReq{
		Name:     &form.Name,
		NickName: &form.NickName,
		Role:     &form.Role,
	}
	res, err := svc.Add(c, &req)
	if err != nil {
		logrus.Error("err", err)
	}
	appG.Response(http.StatusOK, e.SUCCESS, res)
}

type UpdateUserForm struct {
	ID       string `form:"id" valid:"Required;MaxSize(100)"`
	Name     string `form:"name" valid:"Required;MaxSize(100)"`
	NickName string `form:"nickName" valid:"Required;MaxSize(255)"`
	Role     string `form:"role" valid:"Required;MaxSize(255)"`
}

// @Summary Update User
// @Produce  json
// @Accept multipart/form-data
// @Param id path string true "id"
// @Param name formData string true "name"
// @Param nickName formData string true "nick name"
// @Param role formData string true "role"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/users/{id} [put]
func UpdateUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var form UpdateUserForm
	form.ID = com.StrTo(c.Param("id")).String()
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
		logrus.Error("err", err)
		os.Exit(1)
	}
	repo := repository.New(db)
	svc := admin.NewBasicUserService(repo.User)
	req := model.UserReq{
		Name:     &form.Name,
		NickName: &form.NickName,
		Role:     &form.Role,
	}
	res, err := svc.Update(c, form.ID, &req)
	if err != nil {
		logrus.Error("err", err)
	}
	appG.Response(http.StatusOK, e.SUCCESS, res)
}

type DeleteUserForm struct {
	ID string `form:"id" valid:"Required;MaxSize(100)"`
}

// @Summary Delete User
// @Produce  json
// @Param id path string true "id"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /admin/users/{id} [delete]
func DeleteUser(c *gin.Context) {
	appG := app.Gin{C: c}
	var form DeleteUserForm
	form.ID = com.StrTo(c.Param("id")).String()
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
		logrus.Error("err", err)
		os.Exit(1)
	}
	repo := repository.New(db)
	svc := admin.NewBasicUserService(repo.User)
	err = svc.Delete(c, form.ID)
	if err != nil {
		logrus.Error("err", err)
	}
	appG.Response(http.StatusOK, e.SUCCESS, nil)
}
