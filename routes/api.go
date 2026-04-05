// src/routes/api.go
package routes

import (
	"si-baper-backend/controllers"
	"si-baper-backend/middlewares"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
	api := router.Group("/api")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", controllers.Register)
			auth.POST("/login", controllers.Login)
		}

		protected := api.Group("/")
		protected.Use(middlewares.RequireAuth())
		{
			protected.GET("/profile", func(c *gin.Context) {
				userID, _ := c.Get("user_id")
				roleID, _ := c.Get("role_id")

				c.JSON(200, gin.H{
					"message": "Selamat datang di area terproteksi SI-BAPER!",
					"data": gin.H{
						"user_id": userID,
						"role_id": roleID,
					},
				})
			})

			items := protected.Group("/items")
			{
				items.GET("", controllers.GetItems)
				items.POST("", controllers.CreateItem)
				items.PUT("/:id", controllers.UpdateItem)
				items.DELETE("/:id", controllers.DeleteItem)
			}

			categories := protected.Group("/categories")
			{
				categories.GET("", controllers.GetCategories)
				categories.POST("", controllers.CreateCategory)
			}

			requests := protected.Group("/requests")
			{
				requests.POST("", controllers.CreateRequest)
				requests.GET("/my-history", controllers.GetMyRequests)
			}

			approvals := protected.Group("/approvals")
			{
				approvals.GET("/pending", controllers.GetPendingRequests)
				approvals.POST("/process", controllers.ProcessApproval)
			}

			inventory := protected.Group("/inventory")
			{
				inventory.GET("/logs", controllers.GetInventoryLogs)
			}

			dashboard := protected.Group("/dashboard")
			{
				dashboard.GET("/stats", controllers.GetDashboardStats)
			}

			users := protected.Group("/users")
			{
				users.GET("/profile", controllers.GetMyProfile)
				users.PUT("/profile", controllers.UpdateMyProfile)

				users.GET("", controllers.GetAllUsers)
				users.GET("/roles", controllers.GetRoles)
				users.PUT("/:id", controllers.UpdateUser)
				users.DELETE("/:id", controllers.DeleteUser)
			}
		}
	}
}