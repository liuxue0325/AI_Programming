package handler

import (
	"net/http"
	"strconv"
	"videohub/backend/internal/pkg/response"
	"videohub/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	service *service.Service
}

func NewHistoryHandler(service *service.Service) *HistoryHandler {
	return &HistoryHandler{service: service}
}

// GetHistory 获取观看历史
// @Summary 获取观看历史
// @Description 获取观看历史列表
// @Tags history
// @Accept json
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=[]model.WatchHistory}
// @Router /api/history [get]
func (h *HistoryHandler) GetHistory(c *gin.Context) {
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

	historyList, total, err := h.service.Media.GetWatchHistory(c.Request.Context(), page, pageSize)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"items": historyList,
		"total": total,
	})
}

// AddHistory 添加观看历史
// @Summary 添加观看历史
// @Description 记录观看进度
// @Tags history
// @Accept json
// @Produce json
// @Param request body AddHistoryRequest true "历史记录请求"
// @Success 200 {object} response.Response
// @Router /api/history [post]
func (h *HistoryHandler) AddHistory(c *gin.Context) {
	var req struct {
		MediaID   int64   `json:"media_id" binding:"required"`
		EpisodeID *int64  `json:"episode_id"`
		Progress  float64 `json:"progress"`
		Completed bool    `json:"completed"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "invalid request body")
		return
	}

	if err := h.service.Media.AddWatchHistory(c.Request.Context(), req.MediaID, req.EpisodeID, req.Progress, req.Completed); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ClearHistory 清空观看历史
// @Summary 清空观看历史
// @Description 清空所有观看历史
// @Tags history
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/history [delete]
func (h *HistoryHandler) ClearHistory(c *gin.Context) {
	if err := h.service.Media.ClearWatchHistory(c.Request.Context()); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
