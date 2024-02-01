CREATE TABLE IF NOT EXISTS user_balances(
user_id SERIAL PRIMARY KEY ,
balance NUMERIC CHECK (balance >= 0),
created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);