package handler

import (
	"net/http"
	"strconv"
	"videohub/backend/internal/model"
	"videohub/backend/internal/pkg/response"
	"videohub/backend/internal/service"

	"github.com/gin-gonic/gin"
)

type MediaHandler struct {
	service *service.Service
}

func NewMediaHandler(service *service.Service) *MediaHandler {
	return &MediaHandler{service: service}
}

// GetMediaList 获取媒体列表
// @Summary 获取媒体列表
// @Description 获取媒体列表，支持筛选和分页
// @Tags media
// @Accept json
// @Produce json
// @Param type query string false "媒体类型"
// @Param year query int false "年份"
// @Param keyword query string false "搜索关键字"
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Success 200 {object} response.Response{data=service.MediaListResponse}
// @Router /api/media [get]
func (h *MediaHandler) GetMediaList(c *gin.Context) {
	var req service.MediaListRequest

	// 解析查询参数
	if typeStr := c.Query("type"); typeStr != "" {
		mediaType := model.MediaType(typeStr)
		req.Type = &mediaType
	}

	if yearStr := c.Query("year"); yearStr != "" {
		if year, err := strconv.Atoi(yearStr); err == nil {
			req.Year = &year
		}
	}

	if keyword := c.Query("keyword"); keyword != "" {
		req.Keyword = &keyword
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			req.Page = page
		}
	} else {
		req.Page = 1
	}

	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if pageSize, err := strconv.Atoi(pageSizeStr); err == nil && pageSize > 0 {
			req.PageSize = pageSize
		}
	} else {
		req.PageSize = 20
	}

	// 调用服务
	resp, err := h.service.Media.GetMediaList(c.Request.Context(), &req)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, resp)
}

// GetMediaByID 获取媒体详情
// @Summary 获取媒体详情
// @Description 根据ID获取媒体详情
// @Tags media
// @Accept json
// @Produce json
// @Param id path int true "媒体ID"
// @Success 200 {object} response.Response{data=model.Media}
// @Router /api/media/{id} [get]
func (h *MediaHandler) GetMediaByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "invalid media ID")
		return
	}

	media, err := h.service.Media.GetMediaByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrMediaNotFound {
			response.NotFound(c, "media not found")
		} else {
			response.ServerError(c, err.Error())
		}
		return
	}

	response.Success(c, media)
}

// GetMediaPlayInfo 获取媒体播放信息
// @Summary 获取媒体播放信息
// @Description 获取媒体的播放URL和字幕信息
// @Tags media
// @Accept json
// @Produce json
// @Param id path int true "媒体ID"
// @Param episode_id query int false "剧集ID"
// @Success 200 {object} response.Response{data=service.PlayInfo}
// @Router /api/media/{id}/play [get]
func (h *MediaHandler) GetMediaPlayInfo(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "invalid media ID")
		return
	}

	var episodeID *int64
	if episodeIDStr := c.Query("episode_id"); episodeIDStr != "" {
		eid, err := strconv.ParseInt(episodeIDStr, 10, 64)
		if err != nil {
			response.ParamError(c, "invalid episode ID")
			return
		}
		episodeID = &eid
	}

	playInfo, err := h.service.Media.GetPlayInfo(c.Request.Context(), id, episodeID)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, playInfo)
}

// GetEpisodes 获取剧集列表
// @Summary 获取剧集列表
// @Description 获取电视剧的剧集列表
// @Tags media
// @Accept json
// @Produce json
// @Param id path int true "媒体ID"
// @Param season query int false "季数"
// @Success 200 {object} response.Response{data=[]model.Episode}
// @Router /api/media/{id}/episodes [get]
func (h *MediaHandler) GetEpisodes(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "invalid media ID")
		return
	}

	season := 0
	if seasonStr := c.Query("season"); seasonStr != "" {
		if s, err := strconv.Atoi(seasonStr); err == nil && s > 0 {
			season = s
		}
	}

	episodes, err := h.service.Media.GetEpisodes(c.Request.Context(), id, season)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"series_id": id,
		"episodes":  episodes,
	})
}

// GetSubtitles 获取字幕列表
// @Summary 获取字幕列表
// @Description 获取媒体的字幕列表
// @Tags media
// @Accept json
// @Produce json
// @Param id path int true "媒体ID"
// @Success 200 {object} response.Response{data=[]service.Subtitle}
// @Router /api/media/{id}/subtitles [get]
func (h *MediaHandler) GetSubtitles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "invalid media ID")
		return
	}

	subtitles, err := h.service.Media.GetSubtitles(c.Request.Context(), id)
	if err != nil {
		response.ServerError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"media_id":  id,
		"subtitles": subtitles,
	})
}

// DeleteMedia 删除媒体
// @Summary 删除媒体
// @Description 根据ID删除媒体
// @Tags media
// @Accept json
// @Produce json
// @Param id path int true "媒体ID"
// @Success 200 {object} response.Response
// @Router /api/media/{id} [delete]
func (h *MediaHandler) DeleteMedia(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.ParamError(c, "invalid media ID")
		return
	}

	if err := h.service.Media.DeleteMedia(c.Request.Context(), id); err != nil {
		if err == service.ErrMediaNotFound {
			response.NotFound(c, "media not found")
		} else {
			response.ServerError(c, err.Error())
		}
		return
	}

	response.Success(c, nil)
}
