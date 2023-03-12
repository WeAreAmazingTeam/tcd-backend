package transaction

const (
	QueryGetAll = `
		SELECT
			id,
			campaign_id,
			COALESCE(user_id, -1),
			amount,
			status,
			code,
			comment,
			COALESCE(payment_url, '') AS payment_url,
			COALESCE(payment_token, '') AS payment_token,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			transactions
		WHERE
			deleted_at IS NULL
	`

	QueryGetTransactionByID = `
		SELECT
			id,
			campaign_id,
			COALESCE(user_id, -1),
			amount,
			status,
			code,
			comment,
			COALESCE(payment_url, '') AS payment_url,
			COALESCE(payment_token, '') AS payment_token,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			transactions
		WHERE
			deleted_at IS NULL
		AND
			id = ?
		LIMIT
			1
	`

	QueryGetTransactionByCampaignId = `
		SELECT
			id,
			campaign_id,
			user_id,
			(CASE
                WHEN user_id = 0 OR user_id = -1 OR user_id IS NULL THEN "Good Person"
                ELSE COALESCE((SELECT name FROM users WHERE id = user_id), '')
            END) AS user_name,
			amount,
			status,
			code,
			comment,
			COALESCE(payment_url, '') AS payment_url,
			COALESCE(payment_token, '') AS payment_token,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			transactions
		WHERE
			deleted_at IS NULL
		AND
			campaign_id = ?
	`

	QueryGetTransactionByCode = `
		SELECT
			id,
			campaign_id,
			COALESCE(user_id, -1),
			amount,
			status,
			code,
			comment,
			COALESCE(payment_url, '') AS payment_url,
			COALESCE(payment_token, '') AS payment_token,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			transactions
		WHERE
			deleted_at IS NULL
		AND
			code = ?
		LIMIT
			1
	`

	QueryGetTransactionByUserId = `
		SELECT
			id,
			campaign_id,
			COALESCE(user_id, -1),
			amount,
			status,
			code,
			comment,
			COALESCE(payment_url, '') AS payment_url,
			COALESCE(payment_token, '') AS payment_token,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			transactions
		WHERE
			deleted_at IS NULL
		AND
			user_id = ?
	`

	QueryAdminDataTablesTransactions = `
		SELECT
			id,
			campaign_id,
			COALESCE((SELECT title FROM campaigns WHERE id = campaign_id), '') AS campaign_title,
			COALESCE(user_id, -1),
			COALESCE((SELECT name FROM users WHERE id = user_id), '') AS user_name,
			amount,
			status,
			code,
			comment,
			COALESCE(payment_url, '') AS payment_url,
			COALESCE(payment_token, '') AS payment_token,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			transactions
		WHERE
			deleted_at IS NULL
	`

	QueryCountAllAdminDataTablesTransactions = `
		SELECT
			COUNT(id) AS count_id
		FROM
			transactions
		WHERE
			deleted_at IS NULL
	`

	QueryUserDataTablesTransactions = `
		SELECT
			id,
			campaign_id,
			COALESCE((SELECT title FROM campaigns WHERE id = campaign_id), '') AS campaign_title,
			COALESCE(user_id, -1),
			amount,
			status,
			code,
			comment,
			COALESCE(payment_url, '') AS payment_url,
			COALESCE(payment_token, '') AS payment_token,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			transactions
		WHERE
			deleted_at IS NULL
		AND
			user_id = ?
	`

	QueryCountAllUserDataTablesTransactions = `
		SELECT
			COUNT(id) AS count_id
		FROM
			transactions
		WHERE
			deleted_at IS NULL
		AND
			user_id = ?
	`

	QueryGetTotalTransaction = QueryCountAllAdminDataTablesTransactions
)
