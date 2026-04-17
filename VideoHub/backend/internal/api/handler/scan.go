package handler

import (
	"net/http"
	"strconv"
	"videohub/backend/internal/pkg/response"
	"videohub/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type ScanHandler struct {
	service *service.Service
}

func NewScanHandler(service *service.Service) *ScanHandler {
	return &ScanHandler{service: service}
}

// StartScan 启动扫描任务
// @Summary 启动扫描任务
// @Description 启动媒体文件夹扫描
// @Tags scan
// @Accept json
// @Produce json
// @Param request body StartScanRequest true "扫描请求"
// @Success 200 {object} response.Response{data=service.ScanTask}
// @Router /api/scan [post]
func (h *ScanHandler) StartScan(c *gin.Context) {
	var req struct {
		FolderPath string `json:"folder_path" binding:"required"`
		Force      bool   `json:"force"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		response.ParamError(c, "invalid request body")
		return
	}

	task, err := h.service.Scanner.StartScan(c.Request.Context(), req.FolderPath, req.Force)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, task)
}

// GetScanProgress 获取扫描进度
// @Summary 获取扫描进度
// @Description 获取扫描任务的进度
// @Tags scan
// @Accept json
// @Produce json
// @Param task_id path int true "任务ID"
// @Success 200 {object} response.Response{data=service.ScanTask}
// @Router /api/scan/{task_id} [get]
func (h *ScanHandler) GetScanProgress(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		response.ParamError(c, "invalid task ID")
		return
	}

	task, err := h.service.Scanner.GetScanProgress(c.Request.Context(), taskID)
	if err != nil {
		response.NotFound(c, "scan task not found")
		return
	}

	response.Success(c, task)
}

// CancelScan 取消扫描任务
// @Summary 取消扫描任务
// @Description 取消正在进行的扫描任务
// @Tags scan
// @Accept json
// @Produce json
// @Param task_id path int true "任务ID"
// @Success 200 {object} response.Response
// @Router /api/scan/{task_id} [delete]
func (h *ScanHandler) CancelScan(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		response.ParamError(c, "invalid task ID")
		return
	}

	if err := h.service.Scanner.CancelScan(c.Request.Context(), taskID); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
