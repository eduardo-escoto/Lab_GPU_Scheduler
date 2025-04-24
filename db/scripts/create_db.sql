-- Create Database
CREATE DATABASE IF NOT EXISTS gpu_scheduler;


-- Use Database
USE gpu_scheduler;

-- Create Survey Responses Table
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS survey_responses;
SET FOREIGN_KEY_CHECKS = 1;
CREATE TABLE IF NOT EXISTS survey_responses (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL, -- Email to link responses to users
    full_name VARCHAR(255) NOT NULL, -- Full name of the user
    desired_username VARCHAR(255) NOT NULL, -- Desired username
    ssh_key TEXT NOT NULL, -- Public SSH key
    remark TEXT DEFAULT NULL, -- Remark or comment
    user_type ENUM('intern', 'masters', 'phd', 'postdoc', 'faculty', 'visitor') NOT NULL, -- User type
    lab_join_year YEAR DEFAULT NULL, -- Year the user joined the lab
    submitted_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Timestamp for when the survey was submitted
    granted_access_at DATETIME DEFAULT NULL, -- Timestamp for when access was granted
    revoked_access_at DATETIME DEFAULT NULL, -- Timestamp for when access was revoked
    revoke_scheduled_at DATETIME DEFAULT NULL, -- Estimated revocation date
    approving_party VARCHAR(255) DEFAULT NULL, -- Who granted access
    revoking_party VARCHAR(255) DEFAULT NULL, -- Who revoked access
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Timestamp for row creation
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Timestamp for last update
);

-- Create Users Table
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS users;
SET FOREIGN_KEY_CHECKS = 1;
CREATE TABLE IF NOT EXISTS users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE, -- Email address (unique identifier)
    user_name VARCHAR(255) NOT NULL UNIQUE, -- Desired username
    name VARCHAR(255) NOT NULL, -- Full name
    password VARCHAR(255) DEFAULT NULL, -- Password (optional, if needed)
    comment TEXT DEFAULT NULL, -- Remark or comment
    user_type ENUM('intern', 'masters', 'phd', 'postdoc', 'faculty', 'visitor') NOT NULL, -- User type
    lab_join_year YEAR DEFAULT NULL, -- Year the user joined the lab
    access_survey_submitted_at DATETIME DEFAULT NULL, -- Timestamp for the first survey submission
    access_survey_updated_at DATETIME DEFAULT NULL, -- Timestamp for the latest survey submission
    granted_access_at DATETIME DEFAULT NULL, -- When access was granted
    revoked_access_at DATETIME DEFAULT NULL, -- When access was revoked
    revoke_scheduled_at DATETIME DEFAULT NULL, -- Estimated revocation date
    approving_party VARCHAR(255) DEFAULT NULL, -- Who granted access
    revoking_party VARCHAR(255) DEFAULT NULL, -- Who revoked access
    is_admin BOOLEAN DEFAULT FALSE, -- Admin flag
    is_whitelisted BOOLEAN DEFAULT FALSE, -- Whitelist flag
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Timestamp for row creation
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP -- Timestamp for last update
);

-- Create User SSH Keys Table
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS user_ssh_keys;
SET FOREIGN_KEY_CHECKS = 1;
CREATE TABLE IF NOT EXISTS user_ssh_keys (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL, -- Foreign key referencing the users table
    ssh_key TEXT NOT NULL, -- SSH key associated with the user
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Timestamp for when the SSH key was added
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE -- Cascade delete when a user is deleted
);

-- Create Requests Table
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS requests;
SET FOREIGN_KEY_CHECKS = 1;
CREATE TABLE IF NOT EXISTS requests (
    id INT AUTO_INCREMENT PRIMARY KEY,
    user_id INT NOT NULL,
    requested_time INT NOT NULL,
    gpu_size ENUM('small', 'medium', 'large') NOT NULL,
    num_gpus INT NOT NULL,
    priority ENUM('low', 'medium', 'high', 'emergency') NOT NULL,
    server_name VARCHAR(255) DEFAULT NULL,
    status ENUM('scheduled', 'in_progress', 'done', 'cancelled') DEFAULT 'scheduled', -- Updated ENUM values
    start_time DATETIME DEFAULT NULL,
    end_time DATETIME DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create GPUs Table
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS gpus;
SET FOREIGN_KEY_CHECKS = 1;
CREATE TABLE IF NOT EXISTS gpus (
    gpu_uuid CHAR(40) PRIMARY KEY, -- Updated length to 40 characters
    server_name VARCHAR(255) NOT NULL,
    gpu_number INT NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    vram_size_mb INT NOT NULL,
    gpu_serial CHAR(13) DEFAULT NULL UNIQUE, -- Made nullable while keeping the unique constraint
    gpu_bus_id CHAR(16) DEFAULT NULL UNIQUE, -- Made nullable while keeping the unique constraint
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Timestamp for row creation
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Timestamp for last update
    UNIQUE (server_name, gpu_number) -- Retain unique constraint
);

-- Create Request-GPU Assignment Table (Join Table)
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS request_gpu_assignments;
SET FOREIGN_KEY_CHECKS = 1;
CREATE TABLE IF NOT EXISTS request_gpu_assignments (
    request_id INT NOT NULL,
    gpu_uuid CHAR(36) NOT NULL,
    PRIMARY KEY (request_id, gpu_uuid), -- Composite primary key
    FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE CASCADE,
    FOREIGN KEY (gpu_uuid) REFERENCES gpus(gpu_uuid) ON DELETE CASCADE
);

-- Create Whitelist Table
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS whitelist;
SET FOREIGN_KEY_CHECKS = 1;
CREATE TABLE IF NOT EXISTS whitelist (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE
);

-- Create Real-Time Usage Table (Historical Record)
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS real_time_usage;
SET FOREIGN_KEY_CHECKS = 1;
CREATE TABLE IF NOT EXISTS real_time_usage (
    gpu_uuid CHAR(40) NOT NULL, -- Updated to CHAR(40)
    gpu_name VARCHAR(255) NOT NULL, -- Added GPU name (e.g., "NVIDIA Tesla V100")
    server_name VARCHAR(255) NOT NULL,
    gpu_number INT NOT NULL,
    utilization DECIMAL(5,2) NOT NULL, -- GPU utilization percentage (e.g., 75.50 for 75.5%)
    memory_utilization DECIMAL(5,2) NOT NULL, -- Memory utilization percentage (e.g., 60.25 for 60.25%)
    memory_used_mb INT NOT NULL, -- Memory currently in use (in MB)
    memory_available_mb INT NOT NULL, -- Memory available (in MB)
    power_usage_watts DECIMAL(5,2) NOT NULL, -- Power usage in watts (e.g., 150.75 for 150.75W)
    temperature_celsius DECIMAL(5,2) NOT NULL, -- GPU temperature in Celsius (e.g., 65.50 for 65.5°C)
    reported_at DATETIME NOT NULL, -- Manually provided timestamp for historical records
    PRIMARY KEY (gpu_uuid, reported_at), -- Composite primary key to allow multiple records per GPU
    UNIQUE (server_name, gpu_number, reported_at) -- Retain unique constraint with timestamp
);

-- Create GPU Processes Table
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS gpu_processes;
SET FOREIGN_KEY_CHECKS = 1;
CREATE TABLE IF NOT EXISTS gpu_processes (
    id INT AUTO_INCREMENT PRIMARY KEY, -- Unique identifier for each process record
    gpu_uuid CHAR(40) NOT NULL, -- Updated to CHAR(40)
    reported_at DATETIME NOT NULL, -- Timestamp of the GPU usage record (foreign key)
    process_id INT NOT NULL, -- Process ID (PID) of the running process
    process_name VARCHAR(255) NOT NULL, -- Name of the process (e.g., "python")
    user_name VARCHAR(255) NOT NULL, -- User running the process
    gpu_utilization DECIMAL(5,2) NOT NULL, -- Percentage of GPU utilization used by this process
    used_gpu_memory INT NOT NULL, -- Amount of GPU memory used by this process (in MiB)
    FOREIGN KEY (gpu_uuid, reported_at) REFERENCES real_time_usage(gpu_uuid, reported_at) ON DELETE CASCADE,
    FOREIGN KEY (user_name) REFERENCES users(user_name) ON DELETE CASCADE -- Added foreign key
);

-- Create Hourly Historical Usage Table
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS real_time_usage_hourly_historical;
SET FOREIGN_KEY_CHECKS = 1;
CREATE TABLE IF NOT EXISTS real_time_usage_hourly_historical (
    gpu_uuid CHAR(40) NOT NULL, -- Updated to CHAR(40)
    server_name VARCHAR(255) NOT NULL,
    gpu_number INT NOT NULL,
    avg_utilization DECIMAL(5,2) NOT NULL, -- Average GPU utilization
    min_utilization DECIMAL(5,2) NOT NULL, -- Minimum GPU utilization
    max_utilization DECIMAL(5,2) NOT NULL, -- Maximum GPU utilization
    avg_memory_utilization DECIMAL(5,2) NOT NULL, -- Average memory utilization
    min_memory_utilization DECIMAL(5,2) NOT NULL, -- Minimum memory utilization
    max_memory_utilization DECIMAL(5,2) NOT NULL, -- Maximum memory utilization
    avg_memory_used_mb DECIMAL(10,2) NOT NULL, -- Average memory used
    min_memory_used_mb DECIMAL(10,2) NOT NULL, -- Minimum memory used
    max_memory_used_mb DECIMAL(10,2) NOT NULL, -- Maximum memory used
    avg_memory_available_mb DECIMAL(10,2) NOT NULL, -- Average memory available
    min_memory_available_mb DECIMAL(10,2) NOT NULL, -- Minimum memory available
    max_memory_available_mb DECIMAL(10,2) NOT NULL, -- Maximum memory available
    avg_power_usage_watts DECIMAL(5,2) NOT NULL, -- Average power usage
    min_power_usage_watts DECIMAL(5,2) NOT NULL, -- Minimum power usage
    max_power_usage_watts DECIMAL(5,2) NOT NULL, -- Maximum power usage
    max_temperature_celsius DECIMAL(5,2) NOT NULL, -- Maximum temperature
    min_temperature_celsius DECIMAL(5,2) NOT NULL, -- Minimum temperature
    reported_at DATETIME NOT NULL, -- Aggregated hourly timestamp
    PRIMARY KEY (gpu_uuid, reported_at),
    UNIQUE (server_name, gpu_number, reported_at)
);

-- Create Hourly Aggregated GPU Processes Table
SET FOREIGN_KEY_CHECKS = 0;
DROP TABLE IF EXISTS gpu_processes_hourly_historical;
SET FOREIGN_KEY_CHECKS = 1;
CREATE TABLE IF NOT EXISTS gpu_processes_hourly_historical (
    gpu_uuid CHAR(40) NOT NULL, -- Updated to CHAR(40)
    user_name VARCHAR(255) NOT NULL, -- User running the processes
    reported_at DATETIME NOT NULL, -- Aggregated hourly timestamp
    avg_gpu_utilization DECIMAL(5,2) NOT NULL, -- Average GPU utilization by the user
    min_gpu_utilization DECIMAL(5,2) NOT NULL, -- Minimum GPU utilization by the user
    max_gpu_utilization DECIMAL(5,2) NOT NULL, -- Maximum GPU utilization by the user
    avg_used_gpu_memory DECIMAL(10,2) NOT NULL, -- Average GPU memory used by the user (in MiB)
    min_used_gpu_memory DECIMAL(10,2) NOT NULL, -- Minimum GPU memory used by the user (in MiB)
    max_used_gpu_memory DECIMAL(10,2) NOT NULL, -- Maximum GPU memory used by the user (in MiB)
    PRIMARY KEY (gpu_uuid, user_name, reported_at), -- Composite primary key
    FOREIGN KEY (user_name) REFERENCES users(user_name) ON DELETE CASCADE -- Added foreign key
);

