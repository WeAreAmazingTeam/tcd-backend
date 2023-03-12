package user

const (
	QueryGetAllUser = `
		SELECT
			id,
			role,
			name,
			email,
			password,
			e_money,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			users
		WHERE
			deleted_at IS NULL
	`

	QueryAdminDataTablesUsers = `
		SELECT
			id,
			role,
			name,
			email,
			password,
			e_money,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			users
		WHERE
			deleted_at IS NULL
	`

	QueryCountAllAdminDataTablesUsers = `
		SELECT
			COUNT(id) AS count_id
		FROM
			users
		WHERE
			deleted_at IS NULL
	`

	QueryUserRegistered = QueryCountAllAdminDataTablesUsers

	QueryTotalWithdrawalRequest = `
		SELECT
			COUNT(id) AS count_id
		FROM
			user_withdrawal_requests
		WHERE
			deleted_at IS NULL
	`

	QueryDataTablesEMoneyFlow = `
		SELECT
			id,
			status,
			amount,
			note,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			user_emoney_flow
		WHERE
			deleted_at IS NULL
		AND
			user_id = ?
	`

	QueryCountAllDataTablesEMoneyFlow = `
		SELECT
			COUNT(id) AS count_id
		FROM
			user_emoney_flow
		WHERE
			deleted_at IS NULL
		AND
			user_id = ?
	`

	QueryDataTablesWithdrawalRequest = `
		SELECT
			id,
			status,
			amount,
			note,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			user_withdrawal_requests
		WHERE
			deleted_at IS NULL
		AND
			user_id = ?
	`

	QueryCountAllDataTablesWithdrawalRequest = `
		SELECT
			COUNT(id) AS count_id
		FROM
			user_withdrawal_requests
		WHERE
			deleted_at IS NULL
		AND
			user_id = ?
	`

	QueryAdminDataTablesWithdrawalRequest = `
		SELECT
			id,
			user_id,
			COALESCE((SELECT name FROM users WHERE id = user_id), '') AS user_name,
			status,
			amount,
			note,
			created_at,
			created_by,
			updated_at,
			updated_by,
			deleted_at,
			deleted_by
		FROM
			user_withdrawal_requests
		WHERE
			deleted_at IS NULL
	`

	QueryCountAllAdminDataTablesWithdrawalRequest = `
		SELECT
			COUNT(id) AS count_id
		FROM
			user_withdrawal_requests
		WHERE
			deleted_at IS NULL
	`
)
