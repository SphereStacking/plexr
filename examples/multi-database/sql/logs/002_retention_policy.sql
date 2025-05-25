-- Create a function to clean old logs (retention policy)
CREATE OR REPLACE FUNCTION clean_old_logs(retention_days INTEGER DEFAULT 30)
RETURNS void AS $$
BEGIN
    -- Delete application logs older than retention period
    DELETE FROM application_logs 
    WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '1 day' * retention_days;
    
    -- Keep audit logs longer (90 days by default)
    DELETE FROM audit_logs 
    WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '1 day' * (retention_days * 3);
END;
$$ LANGUAGE plpgsql;

-- Create a view for recent errors
CREATE OR REPLACE VIEW recent_errors AS
SELECT 
    id,
    logger_name,
    message,
    stack_trace,
    context,
    created_at
FROM application_logs
WHERE log_level IN ('ERROR', 'FATAL')
  AND created_at > CURRENT_TIMESTAMP - INTERVAL '24 hours'
ORDER BY created_at DESC;