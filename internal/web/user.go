package web

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
	"webook/internal/domain"
	"webook/internal/service"

	regex "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc            *service.UserService
	emailRegexp    *regex.Regexp
	passwordRegexp *regex.Regexp
	userRegexp     *regex.Regexp
	dateRegexp     *regex.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	emailRegexp := regex.MustCompile("^[a-zA-Z0-9_.-]+@[a-zA-Z0-9-]+(\\.[a-zA-Z0-9-]+)*\\.[a-zA-Z0-9]{2,6}$", 0)
	passwordRegexp := regex.MustCompile("^(?=.*\\d)(?=.*[a-zA-Z])(?=.*[^\\da-zA-Z\\s]).{6,18}$", 0)
	userRegexp := regex.MustCompile(`^[a-zA-Z0-9_-]{6,12}$`, 0)
	dateRegexp := regex.MustCompile(`^(?:(?!0000)[0-9]{4}-(?:(?:0[1-9]|1[0-2])-(?:0[1-9]|1[0-9]|2[0-8])|(?:0[13-9]|1[0-2])-(?:29|30)|(?:0[13578]|1[02])-31)|(?:[0-9]{2}(?:0[48]|[2468][048]|[13579][26])|(?:0[48]|[2468][048]|[13579][26])00)-02-29)$`, 0)
	return &UserHandler{
		emailRegexp:    emailRegexp,
		passwordRegexp: passwordRegexp,
		userRegexp:     userRegexp,
		dateRegexp:     dateRegexp,
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

	// 使用 session 进行身份验证
	//session := sessions.Default(ctx)
	//session.Set("userID", user.ID)
	//err = session.Save()
	//if err != nil {
	//	ctx.String(http.StatusInternalServerError, "系统错误")
	//	return
	//}

	// 使用 jwt 进行身份验证
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &UserClaims{
		Uid:       user.ID,
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	})
	tokenStr, err := token.SignedString([]byte("xbcbtlzWUNZfXzmXvnLdpQnoIFRegaUK"))
	if err != nil {
		ctx.String(http.StatusInternalServerError, "系统错误")
		return
	}
	ctx.Header("x-jwt-token", tokenStr)
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
	if errors.Is(err, service.ErrEmailDuplicated) {
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
	type ProfileResp struct {
		Email        string `json:"email"`
		NickName     string `json:"nickName"`
		Avatar       string `json:"avatar"`
		Introduction string `json:"introduction"`
		Birthday     string `json:"birthday"`
	}
	//session := sessions.Default(ctx)
	//id := session.Get("userID").(int64)

	id, _ := ctx.Get("userID")
	user, err := u.svc.Profile(ctx, id.(int64))
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	resp := ProfileResp{
		Email:        user.Email,
		NickName:     user.NickName,
		Avatar:       user.Avatar,
		Introduction: user.Introduction,
		Birthday:     user.Birthday,
	}
	ctx.JSON(http.StatusOK, resp)
}
func (u *UserHandler) Edit(ctx *gin.Context) {
	type EditReq struct {
		NickName     string `json:"nickName"`
		Avatar       string `json:"avatar"`
		Introduction string `json:"introduction"`
		Birthday     string `json:"birthday"`
	}
	var req EditReq
	if err := ctx.Bind(&req); err != nil {
		return
	}
	if len(req.NickName) == 0 {
		ctx.String(http.StatusOK, "昵称不能为空")
		return
	}
	isMatch, err := u.userRegexp.MatchString(req.NickName)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isMatch {
		ctx.String(http.StatusOK, "昵称必须为6-12个字符，且不能包含特殊字符")
		return
	}
	if len([]rune(req.Introduction)) > 255 {
		ctx.String(http.StatusOK, "个人介绍不能多于255个字符")
		return
	}
	if len(req.Birthday) != 0 {
		isMatch, err = u.dateRegexp.MatchString(req.Birthday)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		if !isMatch {
			ctx.String(http.StatusOK, "生日必须为有效日期，日期格式如下：2006-10-10")
			return
		}
	}
	//session := sessions.Default(ctx)
	//id := session.Get("userID"

	id, _ := ctx.Get("userID")
	err = u.svc.Edit(ctx, domain.User{
		ID:           id.(int64),
		NickName:     req.NickName,
		Avatar:       req.Avatar,
		Introduction: req.Introduction,
		Birthday:     req.Birthday,
	})
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	ctx.String(http.StatusOK, "修改成功")
}

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64  `json:"uid"`
	UserAgent string `json:"userAgent"`
}
