package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess      = 0
	CodeParamError   = 1001
	CodeNotFound     = 1002
	CodeForbidden    = 1003
	CodeExists       = 1004
	CodeServerError  = 2001
	CodeDatabaseError = 2002
	CodeFileError    = 2003
	CodeAPIError     = 2004
)

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}

func ParamError(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    CodeParamError,
		Message: message,
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    CodeNotFound,
		Message: message,
	})
}

func ServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeServerError,
		Message: message,
	})
}
