package usecase

type UserUseCase interface {
	Register(username, password string) error
	Login(username, password string) (string, error)
}
