package usecases

type FilterUsecase struct {
	filterRepo FilterRepository
}

func NewFilterUsecase(filterRepo FilterRepository) *FilterUsecase {
	return &FilterUsecase{
		filterRepo: filterRepo,
	}
}
