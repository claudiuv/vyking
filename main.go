package main

import (
	"strconv"
	_ "vyking/docs" // Import your Swagger docs

	"log"
	"net/http"
	"vyking/business" // Import the business layer

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/go-sql-driver/mysql"
)

// @vyking REST API
// @version 1.0
// @description This a RestAPI for an iGaming platform.

func main() {
	/// Initialize the business layer with the database connection string (DSN)
	dbConn, _ := business.GetDatabaseConnection()
	bl, err := business.NewBusinessLayer(dbConn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer bl.Close()
	r := gin.Default()

	// Swagger Route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/api/v1/distributePrize", func(c *gin.Context) { distributePrizeHandler(c, bl) })
	r.GET("/api/v1/rankings", func(c *gin.Context) { getRankingsHandler(c, bl) })

	// Start the server
	r.Run(":8080")
}

// distributePrizeHandler triggers the prize distribution stored procedure
// @Summary Distribute prizes for a tournament
// @Description This endpoint triggers a stored procedure to distribute prizes among players based on their rankings in the specified tournament.
// @Tags Prizes
// @Param tournamentID query int true "ID of the tournament"
// @Success 200 {object} map[string]string "Prizes distributed successfully"
// @Failure 400 {object} map[string]string "Invalid Tournament ID"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/distributePrize [post]
func distributePrizeHandler(c *gin.Context, bl *business.BusinessLayer) {
	tournamentIDStr := c.Query("tournamentID")
	tournamentID, err := strconv.Atoi(tournamentIDStr)
	if err != nil || tournamentID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Tournament ID"})
		return
	}

	err = bl.DistributePrizes(tournamentID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Prizes distributed successfully"})
}

// getRankingsHandler fetches player rankings
// @Summary Fetch player rankings
// @Description This endpoint generates a ranking report based on players' account balances in descending order.
// @Tags Rankings
// @Success 200 {object} map[string]string "List of players with their rankings"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /api/v1/rankings [get]
func getRankingsHandler(c *gin.Context, bl *business.BusinessLayer) {
	rankings, err := bl.GetRankings()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, rankings)
}
