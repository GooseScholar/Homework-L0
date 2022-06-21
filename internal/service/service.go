package service

type Service struct {
	repo Repository
}

//Конструктор сервиса
func New(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}
