-- Create database
CREATE DATABASE xmdb
WITH
ENCODING = 'UTF8'
CONNECTION LIMIT = 200;

-- Connect to newborn db
\c xmdb

-- Create type enum
CREATE TYPE company_type
AS ENUM ('Corporations', 'NonProfit', 'Cooperative', 'Sole Proprietorship');

-- Create table for companies
CREATE TABLE companies (
    id UUID PRIMARY KEY UNIQUE NOT NULL,
    name VARCHAR(15) UNIQUE NOT NULL,
    description VARCHAR(3000),
    employees_amount INT NOT NULL CHECK (employees_amount >= 0),
    registered BOOLEAN NOT NULL,
    type company_type NOT NULL
);

-- Insert some companies
INSERT INTO companies VALUES
    ('2c5165c5-6a66-45fd-a57e-d91726a8b32f', 'Google', 'Google LLC is an American multinational technology company focusing on search engine technology', 30000, 'true', 'Corporations'),
    ('3961b29e-c934-4ab3-b91f-c460e2c86d51', 'Apple', '', 25000, 'true', 'Cooperative'),
    ('7b18f7e6-c9a8-4a04-869e-482ef1c9523d', 'Samsung', 'Samsung is a South Korean multinational manufacturing conglomerate headquartered in Samsung Town, Seoul, South Korea', 15750, 'false', 'Sole Proprietorship');

-- Create index on id
CREATE INDEX idx_id ON companies USING btree(id);
