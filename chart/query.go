package chart

const (
	QueryGetChart = `
		SELECT
			id,
			name,
			content,
			created_at
		FROM
			charts
		WHERE
			name = ?
	`
)
