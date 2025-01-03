package main

import (
	"maker-checker/config"
	messageCtrl "maker-checker/internal/delivery/http/v1/message"
	"maker-checker/internal/middleware"
	messageRepo "maker-checker/internal/repositories/message"
	messageService "maker-checker/internal/services/message"

	userCtrl "maker-checker/internal/delivery/http/v1/user"
	userRepo "maker-checker/internal/repositories/user"
	userService "maker-checker/internal/services/user"

	"github.com/gin-gonic/gin"

	"gorm.io/gorm"
)

func messageRoutes(cfg *config.AppConfig, db *gorm.DB, group *gin.RouterGroup, mw *middleware.Middleware) {
	grp := group.Group("/messages")
	{
		repo := messageRepo.New(db)
		service := messageService.New(repo, cfg)
		ctrl := messageCtrl.New(service)

		grp.POST("", mw.AuthMiddleware("maker", "admin"), ctrl.Create)
		id := grp.Group("/:id")
		{
			id.GET("", ctrl.Get)
			id.PATCH("/approval", mw.AuthMiddleware("checker", "admin"), ctrl.Approval)
		}
	}
}

func userRoutes(cfg *config.AppConfig, db *gorm.DB, group *gin.RouterGroup) {
	grp := group.Group("/users")
	{
		repo := userRepo.New(db)
		service := userService.New(repo, cfg)
		ctrl := userCtrl.New(service)

		grp.POST("", ctrl.Create)
		grp.POST("/login", ctrl.Login)
	}
}
