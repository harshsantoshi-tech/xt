package sql

const GET_USERS_FROM_EMAIL = `SELECT id, username, email, password_hash FROM users WHERE email = ? LIMIT 1`

const UPDATE_PASSWORD_BY_EMAIL = `UPDATE users SET password_hash = ? WHERE email = ?`

const INSERT_USERS = "INSERT INTO users (username , email, password_hash, created_at) VALUES (?,?, ?, ?)"
