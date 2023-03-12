package chart

type Service interface {
	GetChart(chartName, year string) (Chart, error)
}

type service struct {
	repo Repository
}

func NewService(
	repository Repository,
) *service {
	return &service{
		repo: repository,
	}
}
