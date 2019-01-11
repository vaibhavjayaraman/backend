CREATE TABLE user_accounts (
  id SERIAL PRIMARY KEY,
  username VARCHAR(255) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  password_salt TEXT NOT NULL,
  joined TIME,
  created_at TIME,
  updated_at TIME,
  deleted_at TIME
);
