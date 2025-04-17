-- Create Database
CREATE DATABASE gpu_scheduler;

-- Use Database
USE gpu_scheduler;

-- Create Users Table
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    signup_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_admin BOOLEAN DEFAULT FALSE,
    is_whitelisted BOOLEAN DEFAULT FALSE
);

-- Create Requests Table
CREATE TABLE requests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    requested_time INT NOT NULL,
    gpu_size ENUM('small', 'medium', 'large') NOT NULL,
    num_gpus INT NOT NULL,
    priority ENUM('low', 'medium', 'high', 'emergency') NOT NULL,
    server_name VARCHAR(255) DEFAULT NULL,
    status ENUM('pending', 'approved', 'denied', 'completed') DEFAULT 'pending',
    start_time DATETIME DEFAULT NULL,
    end_time DATETIME DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create GPUs Table
CREATE TABLE gpus (
    id CHAR(32) PRIMARY KEY,
    server_name VARCHAR(255) NOT NULL,
    gpu_number INT NOT NULL,
    gpu_size ENUM('small', 'medium', 'large') NOT NULL,
    status ENUM('available', 'in_use', 'maintenance') DEFAULT 'available',
    UNIQUE (server_name, gpu_number)
);

-- Create GPU Usage Table
CREATE TABLE gpu_usage (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    request_id INT NOT NULL,
    start_time DATETIME NOT NULL,
    end_time DATETIME NOT NULL,
    actual_usage_time INT NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE CASCADE
);

-- Create Whitelist Table
CREATE TABLE whitelist (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE
);