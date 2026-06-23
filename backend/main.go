package main

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type TicketPayment struct {
	StationName string  `json:"station_name"`
	Fare        float64 `json:"fare"`
	MomoNumber  string  `json:"momo_number"`
	Provider    string  `json:"provider"`
}

type DriverApplication struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	License      string `json:"license"`
	LicenseClass string `json:"license_class"` // e.g., A, B, C, D
	BirthYear    int    `json:"birth_year"`
	Plate        string `json:"plate"`
	Photo        string `json:"photo"`
	Status       string `json:"status"`
	Reason       string `json:"reason"`
}

type LostItem struct {
	Plate string `json:"plate"`
	Item  string `json:"item"`
	Phone string `json:"phone"`
}

var driversRoster = []DriverApplication{
	{ID: 1, Name: "Kwame Mensah", License: "D1-49202-24", LicenseClass: "C", BirthYear: 1992, Plate: "WR-302-25", Photo: "", Status: "APPROVED", Reason: "Passed initial check"},
}

var lostItemsDatabase = []LostItem{}

func main() {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.POST("/api/tickets/pay", func(c *gin.Context) {
		var req TicketPayment
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "SUCCESS", "details": req})
	})

	// REVISED: Evaluates DVLA criteria rules before approving drivers
	router.POST("/api/drivers/apply", func(c *gin.Context) {
		var req DriverApplication
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		currentYear := time.Now().Year()
		age := currentYear - req.BirthYear

		// Qualification Guard Rails: Must be 18+ and possess Class C or D for commercial passenger vehicles
		if age < 18 {
			req.Status = "REJECTED"
			req.Reason = "Applicant must be at least 18 years old."
		} else if req.LicenseClass != "C" && req.LicenseClass != "D" {
			req.Status = "REJECTED"
			req.Reason = "Commercial driving requires a Class C or Class D license asset."
		} else {
			req.Status = "APPROVED"
			req.Reason = "Qualifications verified successfully by GPRTU standard rules."
		}

		req.ID = time.Now().UnixNano()
		driversRoster = append([]DriverApplication{req}, driversRoster...)
		c.JSON(http.StatusOK, req)
	})

	router.GET("/api/drivers", func(c *gin.Context) {
		c.JSON(http.StatusOK, driversRoster)
	})

	router.POST("/api/lost-found", func(c *gin.Context) {
		var req LostItem
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		lostItemsDatabase = append([]LostItem{req}, lostItemsDatabase...)
		c.JSON(http.StatusOK, gin.H{"status": "SUCCESS"})
	})

	router.GET("/api/lost-found", func(c *gin.Context) {
		c.JSON(http.StatusOK, lostItemsDatabase)
	})

	router.Run(":8080")
}
