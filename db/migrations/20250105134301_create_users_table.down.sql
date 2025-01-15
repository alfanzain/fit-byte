-- Drop indexes
DROP INDEX IF EXISTS idx_users_email_password;
DROP INDEX IF EXISTS idx_users_email;

-- Drop the users table
DROP TABLE IF EXISTS users;
