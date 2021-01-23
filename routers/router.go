package routers

import (
	"jim_evaluate/controller/api/admincont"

	v1 "jim_evaluate/controller/api/v1"
	_ "jim_evaluate/docs"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// InitRouter 初始化路由器
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	//r.Use(middleware.Cors())
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//r.POST("/test", Oauth.GetOauthUrl)
	//注册路由组apiv1
	apiv1 := r.Group("/api/v1")
	{

		// User
		userRoute := apiv1.Group("/user")
		{
			userRoute.POST("/getAllUserList", v1.GetUserList)

			userRoute.POST("/UpdateUserInfo", v1.UpdateUserInfo)
			userRoute.POST("/SubmitReason", admincont.SubmitReason)
			userRoute.POST("/GetOrderInfo", v1.GetOrderInfo)
			userRoute.POST("/GetOrderList", v1.GetOrderList)
			userRoute.POST("/UpdateOrderState", v1.UpdateOrderState)
			userRoute.POST("/SendMsg", v1.SendMsg)
			userRoute.POST("/login", v1.Login)
			userRoute.POST("/GetPaper", v1.GetPaper)
			userRoute.POST("/CommitOrder", v1.CommitOrder)
			userRoute.POST("/CheckReason", admincont.CheckReason)
		}
	}

	return r
}
