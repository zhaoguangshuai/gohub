// Package routes 注册路由
package routes

import (
	controllers "gohub/app/http/controllers/api/v1"
	"gohub/app/http/controllers/api/v1/auth"
	"gohub/app/http/middlewares"
	"gohub/pkg/config"

	"github.com/gin-gonic/gin"
)

// RegisterAPIRoutes 注册 API 相关路由
func RegisterAPIRoutes(r *gin.Engine) {

	// 测试一个 v1 的路由组，我们所有的 v1 版本的路由都将存放到这里
	var v1 *gin.RouterGroup
	if len(config.Get("app.api_domain")) == 0 {
		v1 = r.Group("/api/v1")
	} else {
		v1 = r.Group("/v1")
	}
	// v1 := r.Group("/v1")

	// 全局限流中间件：每小时限流。这里是所有 API （根据 IP）请求加起来。
	// 作为参考 Github API 每小时最多 60 个请求（根据 IP）。
	// 测试时，可以调高一点。
	v1.Use(middlewares.LimitIP("200-H"))

	{
		authGroup := v1.Group("/auth")
		// 限流中间件：每小时限流，作为参考 Github API 每小时最多 60 个请求（根据 IP）
		// 测试时，可以调高一点
		authGroup.Use(middlewares.LimitIP("1000-H"))
		{
			// 登录
			lgc := new(auth.LoginController)
			authGroup.POST("/login/using-phone", middlewares.GuestJWT(), lgc.LoginByPhone)
			authGroup.POST("/login/using-password", middlewares.GuestJWT(), lgc.LoginByPassword)
			authGroup.POST("/login/refresh-token", middlewares.AuthJWT(), lgc.RefreshToken)

			// 重置密码
			pwc := new(auth.PasswordController)
			authGroup.POST("/password-reset/using-email", middlewares.GuestJWT(), pwc.ResetByEmail)
			authGroup.POST("/password-reset/using-phone", middlewares.GuestJWT(), pwc.ResetByPhone)

			// 注册用户
			suc := new(auth.SignupController)
			authGroup.POST("/signup/using-phone", middlewares.GuestJWT(), suc.SignupUsingPhone)
			authGroup.POST("/signup/using-email", middlewares.GuestJWT(), suc.SignupUsingEmail)
			authGroup.POST("/signup/phone/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), suc.IsPhoneExist)
			authGroup.POST("/signup/email/exist", middlewares.GuestJWT(), middlewares.LimitPerRoute("60-H"), suc.IsEmailExist)

			// 发送验证码
			vcc := new(auth.VerifyCodeController)
			authGroup.POST("/verify-codes/phone", middlewares.LimitPerRoute("20-H"), vcc.SendUsingPhone)
			authGroup.POST("/verify-codes/email", middlewares.LimitPerRoute("20-H"), vcc.SendUsingEmail)
			// 图片验证码
			authGroup.POST("/verify-codes/captcha", middlewares.LimitPerRoute("50-H"), vcc.ShowCaptcha)
		}

		uc := new(controllers.UsersController)
		v1.GET("/user", middlewares.AuthJWT(), uc.CurrentUser)
		usersGroup := v1.Group("/users")
		{
			usersGroup.GET("", middlewares.AuthJWT(), uc.Inder)
			usersGroup.PUT("", middlewares.AuthJWT(), uc.UpdateProfile)
			usersGroup.PUT("/email", middlewares.AuthJWT(), uc.UpdateEmail)
			usersGroup.PUT("/phone", middlewares.AuthJWT(), uc.UpdatePhone)
			usersGroup.PUT("/password", middlewares.AuthJWT(), uc.UpdatePassword)
			usersGroup.PUT("/avatar", middlewares.AuthJWT(), uc.UpdateAvatar)
		}

		cgc := new(controllers.CategoriesController)
		cgcGroup := v1.Group("/categories")
		{
			cgcGroup.GET("", cgc.Index)
			cgcGroup.POST("", middlewares.AuthJWT(), cgc.Store)
			cgcGroup.PUT("/:id", middlewares.AuthJWT(), cgc.Update)
			cgcGroup.DELETE("/:id", middlewares.AuthJWT(), cgc.Delete)
		}
		tgc := new(controllers.TopicsController)
		tgcGroup := v1.Group("topics")
		{
			tgcGroup.POST("", middlewares.AuthJWT(), tgc.Store)
			tgcGroup.PUT("/:id", middlewares.AuthJWT(), tgc.Update)
			tgcGroup.DELETE("/:id", middlewares.AuthJWT(), tgc.Delete)
			tgcGroup.GET("", tgc.Index)
			tgcGroup.GET("/:id", tgc.Show)
		}

		lsc := new(controllers.LinksController)
		lscGroup := v1.Group("links")
		{
			lscGroup.GET("", lsc.Index)
		}
	}
}
