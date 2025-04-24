-- Aggregate data older than 3 months to hourly resolution
INSERT INTO real_time_usage_hourly_historical (
    gpu_uuid,
    server_name,
    gpu_number,
    avg_utilization,
    min_utilization,
    max_utilization,
    avg_memory_utilization,
    min_memory_utilization,
    max_memory_utilization,
    avg_memory_used_mb,
    min_memory_used_mb,
    max_memory_used_mb,
    avg_memory_available_mb,
    min_memory_available_mb,
    max_memory_available_mb,
    avg_power_usage_watts,
    min_power_usage_watts,
    max_power_usage_watts,
    max_temperature_celsius,
    min_temperature_celsius,
    reported_at
)
SELECT
    gpu_uuid,
    server_name,
    gpu_number,
    AVG(utilization) AS avg_utilization,
    MIN(utilization) AS min_utilization,
    MAX(utilization) AS max_utilization,
    AVG(memory_utilization) AS avg_memory_utilization,
    MIN(memory_utilization) AS min_memory_utilization,
    MAX(memory_utilization) AS max_memory_utilization,
    AVG(memory_used_mb) AS avg_memory_used_mb,
    MIN(memory_used_mb) AS min_memory_used_mb,
    MAX(memory_used_mb) AS max_memory_used_mb,
    AVG(memory_available_mb) AS avg_memory_available_mb,
    MIN(memory_available_mb) AS min_memory_available_mb,
    MAX(memory_available_mb) AS max_memory_available_mb,
    AVG(power_usage_watts) AS avg_power_usage_watts,
    MIN(power_usage_watts) AS min_power_usage_watts,
    MAX(power_usage_watts) AS max_power_usage_watts,
    MAX(temperature_celsius) AS max_temperature_celsius,
    MIN(temperature_celsius) AS min_temperature_celsius,
    DATE_FORMAT(reported_at, '%Y-%m-%d %H:00:00') AS hourly_time
FROM real_time_usage
WHERE reported_at < NOW() - INTERVAL 3 MONTH
GROUP BY gpu_uuid, server_name, gpu_number, DATE_FORMAT(reported_at, '%Y-%m-%d %H:00:00');

-- Delete high-resolution data older than 3 months
DELETE FROM real_time_usage
WHERE reported_at < NOW() - INTERVAL 3 MONTH;

-- 0 0 * * * mysql -u <user> -p<password> gpu_scheduler < /path/to/cleanup_script.sql