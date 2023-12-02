package web

import (
	"net/http"
	"webook/internal/domain"
	"webook/internal/service"

	regex "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc            *service.UserService
	emailRegexp    *regex.Regexp
	passwordRegexp *regex.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	emailRegexp := regex.MustCompile("^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\\.[a-zA-Z0-9-]+)*\\.[a-zA-Z0-9]{2,6}$", 0)
	passwordRegexp := regex.MustCompile("^(?=.*\\d)(?=.*[a-zA-Z])(?=.*[^\\da-zA-Z\\s]).{6,18}$", 0)
	return &UserHandler{
		emailRegexp:    emailRegexp,
		passwordRegexp: passwordRegexp,
		svc:            svc,
	}
}
func (u *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/users")
	ug.GET("/profile", u.Profile)
	ug.POST("/login", u.Login)
	ug.POST("/edit", u.Edit)
	ug.POST("/signup", u.SignUp)
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "账号/邮箱或密码错误")
		return
	}
	session := sessions.Default(ctx)
	session.Set("userId", user.Id)
	session.Save()
	ctx.String(http.StatusOK, "登录成功")
}
func (u *UserHandler) SignUp(ctx *gin.Context) {

	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	isMatch, err := u.emailRegexp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isMatch {
		ctx.String(http.StatusOK, "邮箱格式错误")
		return
	}
	isMatch, err = u.passwordRegexp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isMatch {
		ctx.String(http.StatusOK, "密码格式错误，至少包含字母、数字、特殊字符，6-18位")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次密码不一致")
		return
	}
	err = u.svc.SignUp(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	if err == service.ErrEmailDuplicated {
		ctx.String(http.StatusOK, "该邮箱已被注册，请重新输入邮箱")
		return
	}
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "注册成功")
}
func (u *UserHandler) Profile(ctx *gin.Context) {

}
func (u *UserHandler) Edit(ctx *gin.Context) {

}
