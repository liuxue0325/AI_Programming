package handler

import (
	"net/http"
	"strconv"
	"videohub/backend/internal/pkg/response"
	"videohub/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	service *service.Service
}

func NewFavoriteHandler(service *service.Service) *FavoriteHandler {
	return &FavoriteHandler{service: service}
}

// GetFavorites 获取收藏列表
// @Summary 获取收藏列表
// @Description 获取收藏的媒体列表
// @Tags favorite
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=[]model.Favorite}
// @Router /api/favorites [get]
func (h *FavoriteHandler) GetFavorites(c *gin.Context) {
	page := 1
	pageSize := 20

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 {
			pageSize = ps
		}
	}

	favoriteList, total, err := h.service.Media.GetFavorites(c.Request.Context(), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"items": favoriteList,
		"total": total,
	})
}

// AddFavorite 添加收藏
// @Summary 添加收藏
// @Description 收藏媒体
// @Tags favorite
// @Accept json
// @Produce json
// @Param request body AddFavoriteRequest true "收藏请求"
// @Success 200 {object} response.Response
// @Router /api/favorites [post]
func (h *FavoriteHandler) AddFavorite(c *gin.Context) {
	var req struct {
		MediaID int64 `json:"media_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "invalid request body")
		return
	}

	if err := h.service.Media.AddFavorite(c.Request.Context(), req.MediaID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// DeleteFavorite 删除收藏
// @Summary 删除收藏
// @Description 取消收藏媒体
// @Tags favorite
// @Accept json
// @Produce json
// @Param media_id path int true "媒体ID"
// @Success 200 {object} response.Response
// @Router /api/favorites/{media_id} [delete]
func (h *FavoriteHandler) DeleteFavorite(c *gin.Context) {
	mediaIDStr := c.Param("media_id")
	mediaID, err := strconv.ParseInt(mediaIDStr, 10, 64)
	if err != nil {
		response.ParamError(c, "invalid media ID")
		return
	}

	if err := h.service.Media.DeleteFavorite(c.Request.Context(), mediaID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
