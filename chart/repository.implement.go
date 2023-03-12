package chart

import (
	"fmt"

	"github.com/WeAreAmazingTeam/tcd-backend/helper"
)

func (repo *repository) GetChart(chartName, year string) (dataChart Chart, err error) {
	additionalQuery := ""

	if year != "" {
		additionalQuery = fmt.Sprintf(" AND created_at LIKE '%%%v-%%'", year)
	}

	additionalQuery += " ORDER BY id DESC LIMIT 1"

	row := repo.DB.Raw(helper.ConvertToInLineQuery(QueryGetChart+additionalQuery), chartName).Row()

	if row.Err() != nil {
		return dataChart, row.Err()
	}

	if err := row.Scan(&dataChart.ID, &dataChart.Name, &dataChart.Content, &dataChart.CreatedAt); err != nil {
		return dataChart, err
	}

	return dataChart, nil
}
