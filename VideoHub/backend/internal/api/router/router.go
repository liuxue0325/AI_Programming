package router

import (
	"videohub/backend/internal/api/handler"
	"videohub/backend/internal/api/middleware"
	"videohub/backend/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func SetupRouter(service *service.Service, logger *zap.Logger) *gin.Engine {
	// 设置Gin模式
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	// 中间件
	r.Use(gin.Recovery())
	r.Use(middleware.Logger(logger))
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// 静态文件
	r.Static("/hls", "./data/hls")
	r.Static("/posters", "./data/posters")
	r.Static("/backdrops", "./data/backdrops")

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 媒体相关
		mediaHandler := handler.NewMediaHandler(service)
		api.GET("/media", mediaHandler.GetMediaList)
		api.GET("/media/:id", mediaHandler.GetMediaByID)
		api.GET("/media/:id/play", mediaHandler.GetMediaPlayInfo)
		api.GET("/media/:id/episodes", mediaHandler.GetEpisodes)
		api.GET("/media/:id/subtitles", mediaHandler.GetSubtitles)
		api.DELETE("/media/:id", mediaHandler.DeleteMedia)

		// 扫描相关
		scanHandler := handler.NewScanHandler(service)
		api.POST("/scan", scanHandler.StartScan)
		api.GET("/scan/:task_id", scanHandler.GetScanProgress)
		api.DELETE("/scan/:task_id", scanHandler.CancelScan)

		// 历史记录相关
		historyHandler := handler.NewHistoryHandler(service)
		api.GET("/history", historyHandler.GetHistory)
		api.POST("/history", historyHandler.AddHistory)
		api.DELETE("/history", historyHandler.ClearHistory)

		// 收藏相关
		favoriteHandler := handler.NewFavoriteHandler(service)
		api.GET("/favorites", favoriteHandler.GetFavorites)
		api.POST("/favorites", favoriteHandler.AddFavorite)
		api.DELETE("/favorites/:media_id", favoriteHandler.DeleteFavorite)

		// 设置相关
		settingHandler := handler.NewSettingHandler(service)
		api.GET("/settings", settingHandler.GetSettings)
		api.PUT("/settings", settingHandler.UpdateSettings)
		api.POST("/settings/reset", settingHandler.ResetSettings)

		// WebSocket
		// api.GET("/ws", websocketHandler.HandleWebSocket)
	}

	return r
}
