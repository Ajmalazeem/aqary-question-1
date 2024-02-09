package db

import (
	"context"
	"database/sql"
	"errors"

	// "fmt"
	"time"

	 "github.com/lib/pq"
)

type User struct {
	ID                int       `json:"id"`
	Name              string    `json:"name"`
	PhoneNumber       string    `json:"phone_number"`
	OTP               string    `json:"otp"`
	OTPExpirationTime *time.Time `json:"otp_expiration_time"`
}

const createUserQuery = `
INSERT INTO users (name, phone_number) VALUES ($1, $2) RETURNING id, name, phone_number, otp, otp_expiration_time;
`

const generateOTPQuery = `
UPDATE users SET otp = $1, otp_expiration_time = $2 WHERE phone_number = $3 RETURNING id, name, phone_number, otp, otp_expiration_time;
`

const verifyOTPQuery = `
SELECT otp, otp_expiration_time FROM users WHERE phone_number = $1;
`

func CreateUser(ctx context.Context, db *sql.DB, name, phoneNumber string) (*User, error) {
	var user User
	err := db.QueryRowContext(ctx, createUserQuery, name, phoneNumber).Scan(
		&user.ID, &user.Name, &user.PhoneNumber, &user.OTP, &user.OTPExpirationTime,
	)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			ErrPhoneNumberExists := errors.New("error 400: phone number alredy exists")
			return nil, ErrPhoneNumberExists
		}
		return nil, err
	}
	return &user, nil
}

func GenerateOTP(ctx context.Context, db *sql.DB, otp, phoneNumber string, expirationTime time.Time) (*User, error) {

	var ErrPhoneNumberNotFound = errors.New("error 404: phone number not found");
	var user User
	err := db.QueryRowContext(ctx, generateOTPQuery, otp, expirationTime, phoneNumber).Scan(
		&user.ID, &user.Name, &user.PhoneNumber, &user.OTP, &user.OTPExpirationTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrPhoneNumberNotFound
		}
		return nil, err
	}
	return &user, nil
}

func VerifyOTP(ctx context.Context, db *sql.DB, phoneNumber string) (*User, error) {
	var user User
	err := db.QueryRowContext(ctx, verifyOTPQuery, phoneNumber).Scan(
		&user.OTP, &user.OTPExpirationTime,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
