package handler

import (
	"WebAnalyzer/internal/dto"
	"github.com/gin-gonic/gin"
)

// analyzePage godoc
// @Summary Analyze a Web Page
// @Schemas
// @Description
// @Tags WebPage
// @Accept  json
// @Produce  json
// @Param input body dto.AnalyzePageReq true "Request Body"
// @Success 200 {object} dto.AnalyzePageRes "Result"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /web-pages/Analyze [post]
func (h *httpHandler) analyzePage() gin.HandlerFunc {
	return func(c *gin.Context) {
		req := &dto.AnalyzePageReq{}
		if err := c.ShouldBindJSON(req); err != nil {
			h.logger.Errorf("[Web Page Handler] %v", err)
			c.JSON(400, gin.H{"error": "error when parsing request data"})
			return
		}

		res, err := h.wpSvc.AnalyzePage(c, req)
		if err != nil {
			h.logger.Errorf("[Web Page Handler] %v", err)
			c.JSON(500, gin.H{"error": "internal server error"})
		}

		c.JSON(200, res)
	}
}
