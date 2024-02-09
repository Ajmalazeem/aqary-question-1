// service.go
package service

import (
	"context"
	"database/sql"
	"log"
	"math/rand"
	"net/http"
	"time"

    "../route"
	"github.com/gin-gonic/gin"
)

// type Service interface {
// 	CreateUser(c *gin.Context)
// 	GetPackageNameDetails(req model.GetRequest) (*model.Model, error)
// 	GetChangeLogDetails(req model.GetRequest) (*[]model.Changelog, error)
// }


type Service struct {
	db *sql.DB
}

func NewService(db *sql.DB) *Service {
	return &Service{db: db}
}

func (s *Service) CreateUser(c *gin.Context) {
	var user db.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	newUser, err := db.CreateUser(ctx, s.db, user.Name, user.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, newUser)
}

func (s *Service) GenerateOTP(c *gin.Context) {
	var input struct {
		PhoneNumber string `json:"phone_number"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	expirationTime := time.Now().Add(1 * time.Minute)
	otp := generateRandomOTP()

	newUser, err := db.GenerateOTP(ctx, s.db, otp, input.PhoneNumber, expirationTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, newUser)
}

func (s *Service) VerifyOTP(c *gin.Context) {
	var input struct {
		PhoneNumber string `json:"phone_number"`
		OTP          string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	user, err := db.VerifyOTP(ctx, s.db, input.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	if user.OTP != input.OTP {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	if user.OTPExpirationTime.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "OTP has expired"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully"})
}

func generateRandomOTP() string {
	// Implement your logic to generate a random 4-digit OTP
	return fmt.Sprintf("%04d", rand.Intn(10000))
}
