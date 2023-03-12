package chart

func (svc *service) GetChart(chartName, year string) (Chart, error) {
	chart, err := svc.repo.GetChart(chartName, year)

	if err != nil {
		return chart, err
	}

	return chart, nil
}
