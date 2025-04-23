-- Aggregate data older than 3 months to hourly resolution
INSERT INTO real_time_usage_hourly_historical (
    gpu_uuid,
    server_name,
    gpu_number,
    utilization,
    memory_utilization,
    memory_used_mb,
    memory_available_mb,
    power_usage_watts,
    temperature_celsius,
    reported_at
)
SELECT
    gpu_uuid,
    server_name,
    gpu_number,
    AVG(utilization) AS avg_utilization,
    AVG(memory_utilization) AS avg_memory_utilization,
    SUM(memory_used_mb) AS total_memory_used,
    SUM(memory_available_mb) AS total_memory_available,
    AVG(power_usage_watts) AS avg_power_usage,
    MAX(temperature_celsius) AS max_temperature,
    DATE_FORMAT(reported_at, '%Y-%m-%d %H:00:00') AS hourly_time
FROM real_time_usage
WHERE reported_at < NOW() - INTERVAL 3 MONTH
GROUP BY gpu_uuid, server_name, gpu_number, DATE_FORMAT(reported_at, '%Y-%m-%d %H:00:00');

-- Delete high-resolution data older than 3 months
DELETE FROM real_time_usage
WHERE reported_at < NOW() - INTERVAL 3 MONTH;

-- 0 0 * * * mysql -u <user> -p<password> gpu_scheduler < /path/to/cleanup_script.sql