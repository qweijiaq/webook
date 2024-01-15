package web

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/qweijiaq/webook/internal/domain"
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"

	"github.com/qweijiaq/webook/internal/service"
)

// UserHandler 其上定义用户相关路由
type UserHandler struct {
	svc         *service.UserService
	emailExp    *regexp.Regexp
	passwordExp *regexp.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	const (
		emailRegexPattern = `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*`
		// passwordRegexPattern 密码包含至少一位数字，字母和特殊字符，且长度 8-16
		passwordRegexPattern = `^(?![0-9a-zA-Z]+$)(?![a-zA-Z!@#$%^&*]+$)(?![0-9!@#$%^&*]+$)[0-9A-Za-z!@#$%^&*]{8,16}$`
	)

	return &UserHandler{
		svc:         svc,
		emailExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
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
		ConfirmPassword string `json:"confirmPassword"`
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
		return
	}

	if req.ConfirmPassword != req.Password {
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
		return
	}

	// 调用一下 svc 的方法
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	if errors.Is(err, service.ErrUserDuplicateEmail) {
		ctx.String(http.StatusOK, "该邮箱已被注册")
		return
	}

	if err != nil {
		ctx.String(http.StatusOK, "系统异常")
		return
	}

	ctx.String(http.StatusOK, "注册成功")

	// 数据库操作

}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, service.ErrInvalidUserOrPassword) {
		ctx.String(http.StatusOK, "邮箱或密码错误")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}

	// 步骤2
	// 登录成功, 设置 session
	sess := sessions.Default(ctx)
	sess.Set("userId", user.Id)
	sess.Save()

	ctx.String(http.StatusOK, "登录成功")
	return
}

func (u *UserHandler) Edit(ctx *gin.Context) {

}

func (u *UserHandler) Profile(ctx *gin.Context) {

}
