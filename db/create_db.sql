-- Create Database
CREATE DATABASE IF NOT EXISTS gpu_scheduler;

-- Use Database
USE gpu_scheduler;

-- Create Users Table
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    signup_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    is_admin BOOLEAN DEFAULT FALSE,
    is_whitelisted BOOLEAN DEFAULT FALSE
);

-- Create Requests Table
CREATE TABLE IF NOT EXISTS requests (
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
CREATE TABLE IF NOT EXISTS gpus (
    id CHAR(32) PRIMARY KEY,
    server_name VARCHAR(255) NOT NULL,
    gpu_number INT NOT NULL,
    manufacturer VARCHAR(255) NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    vram_size_mb INT NOT NULL,
    UNIQUE (server_name, gpu_number)
);

-- Create GPU Usage Table
CREATE TABLE IF NOT EXISTS gpu_usage (
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
CREATE TABLE IF NOT EXISTS whitelist (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE
);

-- Create Real-Time Usage Table
CREATE TABLE IF NOT EXISTS real_time_usage (
    server_name VARCHAR(255) NOT NULL,
    gpu_number INT NOT NULL,
    utilization DECIMAL(5,2) NOT NULL, -- GPU utilization percentage (e.g., 75.50 for 75.5%)
    memory_utilization DECIMAL(5,2) NOT NULL, -- Memory utilization percentage (e.g., 60.25 for 60.25%)
    memory_used_mb INT NOT NULL, -- Memory currently in use (in MB)
    memory_available_mb INT NOT NULL, -- Memory available (in MB)
    power_usage_watts DECIMAL(5,2) NOT NULL, -- Power usage in watts (e.g., 150.75 for 150.75W)
    temperature_celsius DECIMAL(5,2) NOT NULL, -- GPU temperature in Celsius (e.g., 65.50 for 65.5°C)
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (server_name, gpu_number)
);