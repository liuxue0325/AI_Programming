package handler

import (
	"net/http"
	"videohub/backend/internal/pkg/response"
	"videohub/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type SettingHandler struct {
	service *service.Service
}

func NewSettingHandler(service *service.Service) *SettingHandler {
	return &SettingHandler{service: service}
}

// GetSettings 获取系统设置
// @Summary 获取系统设置
// @Description 获取所有系统设置
// @Tags setting
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=service.SystemSettings}
// @Router /api/settings [get]
func (h *SettingHandler) GetSettings(c *gin.Context) {
	settings, err := h.service.Setting.GetSettings(c.Request.Context())
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, settings)
}

// UpdateSettings 更新系统设置
// @Summary 更新系统设置
// @Description 更新系统设置
// @Tags setting
// @Accept json
// @Produce json
// @Param request body map[string]interface{} true "设置更新"
// @Success 200 {object} response.Response
// @Router /api/settings [put]
func (h *SettingHandler) UpdateSettings(c *gin.Context) {
	var settings map[string]interface{}
	if err := c.ShouldBindJSON(&settings); err != nil {
		response.ParamError(c, "invalid request body")
		return
	}

	if err := h.service.Setting.UpdateSettings(c.Request.Context(), settings); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}

// ResetSettings 重置系统设置
// @Summary 重置系统设置
// @Description 重置所有系统设置为默认值
// @Tags setting
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /api/settings/reset [post]
func (h *SettingHandler) ResetSettings(c *gin.Context) {
	if err := h.service.Setting.ResetSettings(c.Request.Context()); err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, nil)
}
