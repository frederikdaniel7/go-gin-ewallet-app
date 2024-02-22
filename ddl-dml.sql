CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);

CREATE TABLE password_tokens(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token VARCHAR NOT NULL,
    expired_at VARCHAR NOT NULL DEFAULT NOW() + INTERVAL '10 minutes', 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE SEQUENCE wallet_serial START WITH 1 INCREMENT BY 1;
CREATE TABLE wallets(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    wallet_number CHAR(10) UNIQUE NOT NULL DEFAULT LPAD(nextval('wallet_serial')::text,10,0),
    balance DECIMAL NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE TYPE TYPE_FUNDS_SOURCE AS ENUM ('Wallet', 'Gacha', 'Bank Transfer', 'Credit Card', 'Pay Later')
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    sender_wallet_id BIGINT NOT NULL,
    recipient_wallet_id BIGINT NOT NULL,
    amount DECIMAL NOT NULL,
    source_of_funds TYPE_FUNDS_SOURCE NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (sender_wallet_id) REFERENCES wallets(id),
    FOREIGN KEY (recipient_wallet_id) REFERENCES wallets(id)
);