CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR NOT NULL,
    name VARCHAR NOT NULL,
    password VARCHAR NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP
);
select * from users;
CREATE TABLE password_tokens(
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    token VARCHAR NOT NULL,
    expired_at TIMESTAMP NOT NULL DEFAULT NOW() + INTERVAL '10 minutes', 
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
select * from wallets;
drop table wallets;
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


CREATE TYPE TYPE_FUNDS_SOURCE AS ENUM ('Wallet', 'Reward', 'Bank Transfer', 'Credit Card', 'Cash');
drop table transactions;
CREATE TABLE transactions (
    id BIGSERIAL PRIMARY KEY,
    sender_wallet_id BIGINT ,
    recipient_wallet_id BIGINT NOT NULL,
    amount DECIMAL NOT NULL,
    source_of_funds TYPE_FUNDS_SOURCE NOT NULL,
    descriptions TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP,
    FOREIGN KEY (sender_wallet_id) REFERENCES wallets(id),
    FOREIGN KEY (recipient_wallet_id) REFERENCES wallets(id)
);

INSERT INTO USERS(email,name, password) VALUES 
('usera@gmail.com','usera','$2a$12$MgRJEJstC6e60U4ePvU6JOZGIPu2JO6RPlTToc9LNPBSyTqEpsSRW'),
('userb@gmail.com','userb','$2a$12$vcfRL.7S/KWuS788J.d.dui3sJGW7C34MKUpFZ1eGcVfTiwkHCg2m'),
('userc@gmail.com','userc','$2a$12$ZOm9ToLuIiRVNxDLdoRJ8.VQz6IAONT99uKnKaQkirnt6gmaLci86'),
('userd@gmail.com','userd','$2a$12$WNQQh1Umciiy/.R3lDRbl.2AS5tzZTLkQk8qybAV6MxvOUyErl3WG'),
('usere@gmail.com','usere','$2a$12$7BzVwlyh.PH1q.xchwcl6erx5dt90zpZQm8h17lM/8mNusqUmRqla');

INSERT INTO wallets(user_id,wallet_number, balance) VALUES
(1,  10000000),
(2,  10000000),
(3,  10000000),
(4,  10000000),
(5,  10000000);

insert into transactions (sender_wallet_id, recipient_wallet_id, amount, source_of_funds, descriptions) values
(1,2,200000,'Wallet','test1'),
(1,2,300000,'Wallet','test2'),
(1,2,400000,'Wallet','test3'),
(1,2,500000,'Wallet','test4'),
(1,2,600000,'Wallet','test5'),
(2,1,600000,'Wallet','test6'),
(2,1,500000,'Wallet','test7'),
(2,1,400000,'Wallet','test8'),
(2,1,300000,'Wallet','test9'),
(2,1,200000,'Wallet','test10'),
(3,4,200000,'Bank Transfer','test1'),
(3,4,300000,'Bank Transfer','test2'),
(3,4,400000,'Bank Transfer','test3'),
(3,4,500000,'Bank Transfer','test4'),
(3,4,600000,'Bank Transfer','test5'),
(4,3,600000,'Credit Card','test6'),
(4,3,500000,'Credit Card','test7'),
(4,3,400000,'Credit Card','test8'),
(4,3,300000,'Credit Card','test9'),
(4,3,200000,'Credit Card','test10');


