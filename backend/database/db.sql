DROP SCHEMA IF EXISTS client_db;
CREATE SCHEMA client_db;
USE client_db;

-- Client Table
CREATE TABLE client (
    client_id VARCHAR(50) PRIMARY KEY,
    first_name CHAR(50),
    last_name CHAR(50),
    dob DATE,
    gender VARCHAR(10),
    email VARCHAR(100) UNIQUE,
    phone VARCHAR(15),
    address VARCHAR(100),
    city VARCHAR(50),
    state VARCHAR(50),
    country VARCHAR(50),
    postal_code VARCHAR(10)
);

-- Account Table
CREATE TABLE account (
    account_id int PRIMARY KEY AUTO_INCREMENT,
    client_id VARCHAR(50),
    account_type VARCHAR(50),
    account_status VARCHAR(50),
    opening_date DATE,
    initial_deposit FLOAT,
    currency VARCHAR(50),
    branch_id VARCHAR(50),
    FOREIGN KEY (client_id) REFERENCES client(client_id)
);

-- Transaction Table
CREATE TABLE transaction (
    transaction_id int PRIMARY KEY AUTO_INCREMENT,
    client_id VARCHAR(50),
    transaction VARCHAR(100),
    amount FLOAT,
    date DATE,
    status VARCHAR(50),
    FOREIGN KEY (client_id) REFERENCES client(client_id)
);

-- Log Table
CREATE TABLE log (
    log_id int PRIMARY KEY AUTO_INCREMENT,
    crud VARCHAR(50),
    attribute_name VARCHAR(100),
    before_value VARCHAR(100),
    after_value VARCHAR(100),
    agent_id VARCHAR(50),
    client_id VARCHAR(50),
    datetime DATETIME,
    FOREIGN KEY (client_id) REFERENCES client(client_id)
);

DROP SCHEMA IF EXISTS user_db;
CREATE SCHEMA user_db;
USE user_db;

-- User Table
CREATE TABLE user (
    id VARCHAR(50) PRIMARY KEY,
    first_name CHAR(50),
    last_name CHAR(50),
    email VARCHAR(100) UNIQUE,
    role VARCHAR(50)
);
