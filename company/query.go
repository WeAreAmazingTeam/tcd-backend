package company

const (
	QueryAdminDataTablesCompanyCashFlow = `
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
			company_cash_flow
		WHERE
			deleted_at IS NULL
	`

	QueryCountAllAdminDataTablesCompanyCashFlow = `
		SELECT
			COUNT(id) AS count_id
		FROM
			company_cash_flow
		WHERE
			deleted_at IS NULL
	`
)
