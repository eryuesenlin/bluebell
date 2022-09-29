package routers

import (
	"bluebell/controller"
	"bluebell/logger"
	"bluebell/middlewares"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Setup() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))

	r.LoadHTMLFiles("./templates/index.html")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})
	// 路由组
	v1 := r.Group("/api/v1")
	// 注册业务路由
	v1.POST("/signup", controller.SignUpHandler)
	// 登陆业务路由
	v1.POST("/login", controller.LoginHandler)
	// 应用JWT认证中间件
	v1.Use(middlewares.JWTAuthMiddleware())
	// 根据时间或分数获取帖子列表
	{
		// 社区列表
		v1.GET("/community", controller.CommunityHandler)
		// 根据社区ID获取社区详情
		v1.GET("/community/:id", controller.CommunityDetailHandler)

		// 发帖子
		v1.POST("/post", controller.CreatePostHandler)
		// 根据帖子ID获取详情
		v1.GET("/post/:id", controller.GetPostDetailHandler)
		// 帖子列表
		v1.GET("/posts/", controller.GetPostListHandler)

		// 投票
		v1.POST("/vote", controller.PostVoteHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "404",
		})
	})
	return r
}
