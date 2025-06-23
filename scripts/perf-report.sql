SELECT
	total_time/1000000000 as total_time_ms,
	total_calls,
	total_time/total_calls/1000000000 as avg_time_per_call_ms,
	object_type,
	object_name,
	event_name,
	sql_text
FROM (
	SELECT
		object_type,
		object_name,
		sql_text,
		event_name,
		SUM(TIMER_WAIT) AS total_time,
		COUNT(*) AS total_calls
	FROM
		performance_schema.events_statements_history_long
	GROUP BY sql_text, event_name, object_type, object_name
) t0
ORDER BY total_time DESC
LIMIT 100;
