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

	// --- 1. Auth Routes ---
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.Register(db))
		auth.POST("/login", controllers.Login(db))
	}

	// --- 2. Public API (ไม่ต้อง Login) ---
	apiPublic := r.Group("/api")
	{
		apiPublic.GET("/activities", controllers.ListActivities(db))
		apiPublic.GET("/activities/:id", controllers.GetActivityByID(db))
		apiPublic.GET("/activities/search", controllers.SearchAndFilterActivities(db)) // เพิ่มกลับมา
		apiPublic.GET("/activities/:id/stats", controllers.GetActivityStats(db))       // เพิ่มกลับมา

		apiPublic.GET("/master-goals", controllers.GetActivityMasterGoals(db))
		apiPublic.GET("/master-categories", controllers.GetActivityMasterCategories(db))
	}

	// --- 3. Private API (ต้องเป็น Member หรือ Admin) ---
	apiPrivate := r.Group("/api", middleware.AuthMiddleware("member", "admin"))
	{
		// Profile Management (แก้ปัญหา 404 ที่คุณเจอ)
		apiPrivate.GET("/profile", controllers.GetProfile(db))
		apiPrivate.PUT("/profile", controllers.UpdateProfile(db))               // เพิ่มกลับมา (สำคัญ!)
		apiPrivate.DELETE("/profile/image", controllers.DeleteProfileImage(db)) // เพิ่มกลับมา

		// Favorites
		apiPrivate.POST("/activities/:id/favorite", controllers.ToggleFavorite(db))
		apiPrivate.GET("/favorites", controllers.ListFavorites(db)) // เพิ่มกลับมา

		// Read History
		apiPrivate.POST("/activities/:id/read", controllers.RecordReadHistory(db)) // เพิ่มกลับมา
		apiPrivate.GET("/read-history", controllers.ListReadHistory(db))           // เพิ่มกลับมา
	}

	// --- 4. Admin API (ต้องเป็น Admin เท่านั้น) ---
	admin := r.Group("/admin", middleware.AuthMiddleware("admin"))
	{
		admin.GET("/users", controllers.ListAllUsers(db))

		// Activity Management
		admin.POST("/activities", controllers.CreateActivity(db))
		admin.PUT("/activities/:id", controllers.UpdateActivity(db))    // เพิ่มกลับมา
		admin.DELETE("/activities/:id", controllers.DeleteActivity(db)) // เพิ่มกลับมา

		// User/Role Management
		admin.POST("/roles", controllers.AdminCreateUser(db))       // เพิ่มกลับมา (หมายถึงสร้าง User โดย Admin)
		admin.DELETE("/roles/:id", controllers.AdminDeleteUser(db)) // เพิ่มกลับมา
	}

	return r
}

func getAllowedOrigins() []string {
	if origins := os.Getenv("CORS_ORIGINS"); origins != "" {
		return strings.Split(origins, ",")
	}
	return []string{"http://localhost:4200"}
}
