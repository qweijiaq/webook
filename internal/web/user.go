package web

import (
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserHandler 其上定义用户相关路由
type UserHandler struct {
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler() *UserHandler {
	const (
		emailRegexPattern    = "邮箱校验正则表达式"
		passwordRegexPattern = "密码校验正则表达式"
	)

	emailExp := regexp.MustCompile(emailRegexPattern, regexp.None)
	passwordExp := regexp.MustCompile(emailRegexPattern, regexp.None)

	return &UserHandler{
		emailExp:    emailExp,
		passwordExp: passwordExp,
	}
}

func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")

	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.PUT("/edit", u.Edit)
	ug.GET("/profile", u.Profile)
}

// Signup 注册用户
func (u *UserHandler) Signup(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		confirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	// Bind 方法会根据 Content-Type 来解析数据到 req 中
	// 如果解析失败，会直接写回一个 400 Bad Request 错误
	if err := ctx.Bind(&req); err != nil {
		return
	}

	ok, err := u.emailExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "请输入一个合法的邮箱")
	}

	if req.confirmPassword != req.Password {
		ctx.String(http.StatusOK, "两次输入的密码不一致")
		return
	}

	ok, err = u.passwordExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	if !ok {
		ctx.String(http.StatusOK, "密码必须大于8位，包含数字、字母和特殊字符")
	}

	ctx.String(http.StatusOK, "注册成功")

	// 数据库操作

}

func (u *UserHandler) Login(ctx *gin.Context) {

}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
