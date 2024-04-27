package countComics

type ComixRepo interface {
	Size() (int, error)
}

type UseCase struct {
	comixRepo ComixRepo
}

func New(comixRepo ComixRepo) *UseCase {
	return &UseCase{
		comixRepo: comixRepo,
	}
}

func (u *UseCase) Run() (int, error) {
	size, err := u.comixRepo.Size()
	return size, err
}
