-- Create the database
CREATE DATABASE golang_project_db;

-- Switch to the created database to create tables, etc.
-- Instead of \c, you can directly create the table in the target database
-- This script is executed with user privileges set in the environment

-- Create the table
CREATE TABLE IF NOT EXISTS golang_project_table (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL
);
