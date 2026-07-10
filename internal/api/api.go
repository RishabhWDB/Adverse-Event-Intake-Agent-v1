package api

import (
	"net/http"

	"github.com/RishabhWDB/Adverse-Event-Intake-Agent-v1/internal/store"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func SetupRouter(pool *pgxpool.Pool) *gin.Engine {
	r := gin.Default()

	// GET /cases — fetch all cases
	r.GET("/cases", func(c *gin.Context) {
		cases, err := store.GetAllCases(pool)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, cases)
	})

	r.GET("/cases/:id", func(c *gin.Context) {
		caseID := c.Param("id")
		caseDetail, err := store.GetCaseByID(pool, caseID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "case not found"})
			return
		}
		c.JSON(http.StatusOK, caseDetail)
	})

	// PATCH /cases/:id/status — update a case's status
	r.PATCH("/cases/:id/status", func(c *gin.Context) {
		caseID := c.Param("id")

		var body struct {
			Status string `json:"status"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
			return
		}

		validStatuses := map[string]bool{"pending": true, "reviewed": true, "closed": true}
		if !validStatuses[body.Status] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "status must be pending, reviewed, or closed"})
			return
		}

		if err := store.UpdateCaseStatus(pool, caseID, body.Status); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"case_id": caseID, "status": body.Status})
	})

	r.Static("/static", "./static")
	r.GET("/", func(c *gin.Context) {
		c.File("./static/index.html")
	})

	return r
}
