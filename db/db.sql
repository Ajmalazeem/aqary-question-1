-- CreateUser
INSERT INTO users (name, phone_number) VALUES ($1, $2) RETURNING *;

--CheckPhoneNumberExists
SELECT EXISTS (SELECT 1 FROM users WHERE phone_number = $1);

-- UpdateOTP
UPDATE users SET otp = $2, otp_expiration_time = $3 WHERE phone_number = $1 RETURNING *;

--name: GetUserByPhoneNumber
SELECT * FROM users WHERE phone_number = $1;


-- Create the "users" table
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) UNIQUE NOT NULL,
    otp VARCHAR(6) DEFAULT '',
    otp_expiration_time TIMESTAMP
);