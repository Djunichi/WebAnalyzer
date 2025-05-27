package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// getAll godoc
// @Summary Get All Previous analyses
// @Schemas
// @Description
// @Tags WebPage
// @Accept  json
// @Produce  json
// @Success 200 {object} dto.GetAllAnalysesRes "Result"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /api/v1/analyses/all [get]
func (h *httpHandler) getAll() gin.HandlerFunc {
	return func(c *gin.Context) {
		res, err := h.aSvc.GetAllAnalyses(c)
		if err != nil {
			h.logger.Errorf("[Web Page Handler] %v", err)
			c.JSON(500, gin.H{"error": "internal server error"})
		}

		c.JSON(200, res)
	}
}

// getById godoc
// @Summary Get Previous analysis by ID
// @Schemas
// @Description
// @Tags WebPage
// @Accept  json
// @Produce  json
// @Param analysis-id query string true "analysis UUID"
// @Success 200 {object} dto.AnalyzePageRes "Result"
// @Failure 400 {object} error "Bad Request"
// @Failure 500 {object} error "Internal Server Error"
// @Router /api/v1/analyses/by-id [get]
func (h *httpHandler) getById() gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := uuid.Parse(c.Query("analysis-id"))
		if err != nil {
			c.JSON(400, gin.H{"error": "analysis-id must be a valid uuid"})
			return
		}

		res, err := h.aSvc.GetAnalysisById(c, id)
		if err != nil {
			h.logger.Errorf("[Web Page Handler] %v", err)
			c.JSON(500, gin.H{"error": "internal server error"})
		}

		c.JSON(200, res)
	}
}
