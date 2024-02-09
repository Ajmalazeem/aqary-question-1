-- db.sql
-- name: CreateUser
INSERT INTO users (name, phone_number) VALUES ($1, $2) RETURNING *;

-- name: CheckPhoneNumberExists
SELECT EXISTS (SELECT 1 FROM users WHERE phone_number = $1);

-- name: UpdateOTP
UPDATE users SET otp = $2, otp_expiration_time = $3 WHERE phone_number = $1 RETURNING *;

-- name: GetUserByPhoneNumber
SELECT * FROM users WHERE phone_number = $1;
