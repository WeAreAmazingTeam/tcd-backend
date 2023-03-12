package logs

const (
	QueryAdminDataTablesActivityLogs = `
		SELECT
			id,
			content,
			user_agent,
			ip_address,
			created_at,
			created_by,
			deleted_at,
			deleted_by
		FROM
			activity_logs
		WHERE
			deleted_at IS NULL
	`

	QueryCountAllAdminDataTablesActivityLogs = `
		SELECT
			COUNT(id) AS count_id
		FROM
			activity_logs
		WHERE
			deleted_at IS NULL
	`
)
