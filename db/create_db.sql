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
    status ENUM('scheduled', 'in_progress', 'done', 'cancelled') DEFAULT 'scheduled', -- Updated ENUM values
    start_time DATETIME DEFAULT NULL,
    end_time DATETIME DEFAULT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Create GPUs Table
CREATE TABLE IF NOT EXISTS gpus (
    gpu_uuid CHAR(36) PRIMARY KEY, -- Using GPU_UUID as the primary key
    server_name VARCHAR(255) NOT NULL,
    gpu_number INT NOT NULL,
    manufacturer VARCHAR(255) NOT NULL,
    model_name VARCHAR(255) NOT NULL,
    vram_size_mb INT NOT NULL,
    gpu_serial VARCHAR(255) NOT NULL UNIQUE, -- New column for GPU serial, unique
    gpu_bus_id VARCHAR(255) NOT NULL UNIQUE, -- New column for GPU bus ID, unique
    UNIQUE (server_name, gpu_number) -- Retain unique constraint
);

-- Create Request-GPU Assignment Table (Join Table)
CREATE TABLE IF NOT EXISTS request_gpu_assignments (
    request_id INT NOT NULL,
    gpu_uuid CHAR(36) NOT NULL,
    PRIMARY KEY (request_id, gpu_uuid), -- Composite primary key
    FOREIGN KEY (request_id) REFERENCES requests(id) ON DELETE CASCADE,
    FOREIGN KEY (gpu_uuid) REFERENCES gpus(gpu_uuid) ON DELETE CASCADE
);

-- Create Whitelist Table
CREATE TABLE IF NOT EXISTS whitelist (
    id INT AUTO_INCREMENT PRIMARY KEY,
    email VARCHAR(255) NOT NULL UNIQUE
);

-- Create Real-Time Usage Table (Historical Record)
CREATE TABLE IF NOT EXISTS real_time_usage (
    gpu_uuid CHAR(36) NOT NULL, -- Added GPU_UUID
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
CREATE TABLE IF NOT EXISTS gpu_processes (
    id INT AUTO_INCREMENT PRIMARY KEY, -- Unique identifier for each process record
    gpu_uuid CHAR(36) NOT NULL, -- GPU UUID (foreign key)
    reported_at DATETIME NOT NULL, -- Timestamp of the GPU usage record (foreign key)
    process_id INT NOT NULL, -- Process ID (PID) of the running process
    process_name VARCHAR(255) NOT NULL, -- Name of the process (e.g., "python")
    user_name VARCHAR(255) NOT NULL, -- User running the process
    gpu_utilization DECIMAL(5,2) NOT NULL, -- Percentage of GPU utilization used by this process
    used_gpu_memory INT NOT NULL, -- Amount of GPU memory used by this process (in MiB)
    FOREIGN KEY (gpu_uuid, reported_at) REFERENCES real_time_usage(gpu_uuid, reported_at) ON DELETE CASCADE
);

-- Create Hourly Historical Usage Table
CREATE TABLE IF NOT EXISTS real_time_usage_hourly_historical (
    gpu_uuid CHAR(36) NOT NULL,
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
CREATE TABLE IF NOT EXISTS gpu_processes_hourly_historical (
    gpu_uuid CHAR(36) NOT NULL, -- GPU UUID
    user_name VARCHAR(255) NOT NULL, -- User running the processes
    reported_at DATETIME NOT NULL, -- Aggregated hourly timestamp
    avg_gpu_utilization DECIMAL(5,2) NOT NULL, -- Average GPU utilization by the user
    min_gpu_utilization DECIMAL(5,2) NOT NULL, -- Minimum GPU utilization by the user
    max_gpu_utilization DECIMAL(5,2) NOT NULL, -- Maximum GPU utilization by the user
    avg_used_gpu_memory DECIMAL(10,2) NOT NULL, -- Average GPU memory used by the user (in MiB)
    min_used_gpu_memory DECIMAL(10,2) NOT NULL, -- Minimum GPU memory used by the user (in MiB)
    max_used_gpu_memory DECIMAL(10,2) NOT NULL, -- Maximum GPU memory used by the user (in MiB)
    PRIMARY KEY (gpu_uuid, user_name, reported_at) -- Composite primary key
);