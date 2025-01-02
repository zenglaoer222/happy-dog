package routes

import (
	"github.com/gin-gonic/gin"
	v1 "happy-dog/api/v1"
	"happy-dog/middleware"
	"happy-dog/utils"
)

func InitRouter() {
	gin.SetMode(utils.AppMode)
	r := gin.New()
	r.Use(middleware.Cors())

	private := r.Group("api/v1")
	private.Use(middleware.JwtToken())
	{
		// 顾客模块
		private.GET("customer", v1.GetCustomer)
		private.DELETE("customer/:id", v1.DeleteCustomer)
		// 商家模块

		private.DELETE("shop/:id", v1.DeleteShop)
	}

	protected_shop := r.Group("api/v1")
	protected_shop.Use(middleware.JwtToken())
	{

		// 商品模块
		protected_shop.POST("shop/add", v1.AddShop)

		protected_shop.PUT("shop/:id", v1.EditShop)

		// 订单模块

		// 商品模块
		protected_shop.POST("product/add", v1.CreateProduct)

		// 商品列表模块

	}

	protected_customer := r.Group("api/v1")
	protected_customer.Use(middleware.JwtToken())
	{
		// 顾客模块
		protected_customer.PUT("customer/:id", v1.EditCustomer)

		protected_customer.POST("order/add", v1.CreateOrder)
		protected_customer.GET("order", v1.GetOrder)
		protected_customer.POST("customer/update", v1.Upload)

		// 余额模块
		protected_customer.POST("wallet/add", v1.CreateWallet)
		protected_customer.GET("wallet/balance", v1.InquireBalance)

		// 获取商家
		protected_customer.GET("shop", v1.GetShop)

		// 好友模块
		protected_customer.POST("friends/add", v1.CreateFriends)
		protected_customer.GET("friends", v1.GetFriends)
		protected_customer.GET("friends/wait", v1.GetFriendsWait)
		protected_customer.POST("friends/accept", v1.AcceptFriends)
		protected_customer.GET("friends/search", v1.SearchFriends)

	}
	public := r.Group("api/v1")
	{
		public.POST("customer/add", v1.AddCustomer)
		public.POST("customer/login", v1.Login)
		public.POST("manager/login", v1.ManagerLogin)
		//public.GET("message", v1.ReceiveMsg)

		// 消息模块
		public.GET("message/add", v1.ConCreate)
		public.GET("message/delete", v1.ConDelete)
	}

	r.Run(utils.HttpPort)
}
