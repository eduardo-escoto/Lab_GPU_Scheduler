-- Aggregate data older than 3 months to hourly resolution
INSERT INTO gpu_processes_hourly_historical (
    gpu_uuid,
    user_name,
    reported_at,
    avg_gpu_utilization,
    min_gpu_utilization,
    max_gpu_utilization,
    avg_used_gpu_memory,
    min_used_gpu_memory,
    max_used_gpu_memory
)
SELECT
    gpu_uuid,
    user_name,
    DATE_FORMAT(reported_at, '%Y-%m-%d %H:00:00') AS hourly_time,
    AVG(gpu_utilization) AS avg_gpu_utilization,
    MIN(gpu_utilization) AS min_gpu_utilization,
    MAX(gpu_utilization) AS max_gpu_utilization,
    AVG(used_gpu_memory) AS avg_used_gpu_memory,
    MIN(used_gpu_memory) AS min_used_gpu_memory,
    MAX(used_gpu_memory) AS max_used_gpu_memory
FROM gpu_processes
WHERE reported_at < NOW() - INTERVAL 3 MONTH
GROUP BY gpu_uuid, user_name, DATE_FORMAT(reported_at, '%Y-%m-%d %H:00:00');

-- Delete high-resolution data older than 3 months
DELETE FROM gpu_processes
WHERE reported_at < NOW() - INTERVAL 3 MONTH;

--0 0 * * * mysql -u <user> -p<password> gpu_scheduler < /path/to/cleanup_gpu_processes.sql