-- Create the 'main' user with a password
CREATE USER "main" WITH PASSWORD 'mainpassword' SUPERUSER; 

-- Create the 'minimal-user' with a password
CREATE USER "minimal-user" WITH PASSWORD 'minimalpassword'; 

-- Create the 'fake-movies-db' database
CREATE DATABASE "fake-movies-db";

-- Grant all privileges on the database to the 'main' user
GRANT ALL PRIVILEGES ON DATABASE "fake-movies-db" TO "main"; 

-- It's good practice to also grant the CONNECT privilege
GRANT CONNECT ON DATABASE "fake-movies-db" TO "main";