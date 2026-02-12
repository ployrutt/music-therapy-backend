// package router

// import (
// 	"os"
// 	"project-backend/controllers"
// 	"project-backend/middleware"
// 	"strings"
// 	"time"

// 	"github.com/gin-contrib/cors"
// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// func SetupRouter(db *gorm.DB) *gin.Engine {
// 	r := gin.Default()

// 	// แก้ไขส่วนนี้เพื่อให้ดึงค่า CORS_ORIGINS ที่เราตั้งใน Dokploy มาใช้
// 	r.Use(cors.New(cors.Config{
// 		AllowOrigins:     getAllowedOrigins(),
// 		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
// 		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
// 		AllowCredentials: true,
// 		MaxAge:           12 * time.Hour,
// 	}))

// 	auth := r.Group("/auth")
// 	{
// 		auth.POST("/register", controllers.Register(db))
// 		auth.POST("/login", controllers.Login(db))
// 	}

// 	publicApi := r.Group("/api")
// 	{
// 		publicApi.GET("/activities", controllers.ListActivities(db))
// 		publicApi.GET("/activities/:id", controllers.GetActivityByID(db))

// 		publicApi.GET("/master-goals", controllers.GetActivityMasterGoals(db))
// 		publicApi.GET("/master-categories", controllers.GetActivityMasterCategories(db))

// 		publicApi.GET("/activities/search", controllers.SearchAndFilterActivities(db))

// 		publicApi.GET("/activities/:id/stats", controllers.GetActivityStats(db))
// 	}

// 	api := r.Group("/api", middleware.AuthMiddleware("member", "admin"))
// 	{
// 		api.GET("/profile", controllers.GetProfile(db))
// 		api.PUT("/profile", controllers.UpdateProfile(db))
// 		api.DELETE("/profile/image", controllers.DeleteProfileImage(db))

// 		api.POST("/activities/:id/favorite", controllers.ToggleFavorite(db))
// 		api.GET("/favorites", controllers.ListFavorites(db))

// 		api.POST("/activities/:id/read", controllers.RecordReadHistory(db))
// 		api.GET("/read-history", controllers.ListReadHistory(db))

// 	}

// 	admin := r.Group("/admin", middleware.AuthMiddleware("admin"))
// 	{
// 		admin.GET("/users", controllers.ListAllUsers(db))

// 		admin.POST("/activities", controllers.CreateActivity(db))
// 		admin.PUT("/activities/:id", controllers.UpdateActivity(db))
// 		admin.DELETE("/activities/:id", controllers.DeleteActivity(db))

// 		admin.POST("/roles", controllers.AdminCreateUser(db))
// 		admin.DELETE("/roles/:id", controllers.AdminDeleteUser(db))
// 	}

// 	return r
// }

//	func getAllowedOrigins() []string {
//		if origins := os.Getenv("CORS_ORIGINS"); origins != "" {
//			return strings.Split(origins, ",")
//		}
//		return []string{"http://localhost:4200"}
//	}
package router

import (
	"os"
	"project-backend/controllers"
	"project-backend/middleware"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     getAllowedOrigins(),
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// --- Routes ---
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register(db))
		auth.POST("/login", controllers.Login(db))
	}

	apiPublic := r.Group("/api")
	{
		apiPublic.GET("/activities", controllers.ListActivities(db))
		apiPublic.GET("/activities/:id", controllers.GetActivityByID(db))
		apiPublic.GET("/master-goals", controllers.GetActivityMasterGoals(db))
		apiPublic.GET("/master-categories", controllers.GetActivityMasterCategories(db))

	}

	apiPrivate := r.Group("/api", middleware.AuthMiddleware("member", "admin"))
	{
		apiPrivate.GET("/profile", controllers.GetProfile(db))
		apiPrivate.POST("/activities/:id/favorite", controllers.ToggleFavorite(db))
	}

	admin := r.Group("/admin", middleware.AuthMiddleware("admin"))
	{
		admin.GET("/users", controllers.ListAllUsers(db))
		admin.POST("/activities", controllers.CreateActivity(db))
	}

	return r
}

func getAllowedOrigins() []string {
	if origins := os.Getenv("CORS_ORIGINS"); origins != "" {
		return strings.Split(origins, ",")
	}
	return []string{"http://localhost:4200"}
}
