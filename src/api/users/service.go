package users

var (
	UserService UserServiceInterface = &userService{}
)

type UserServiceInterface interface {
	Create(user *User) (*UserDto, error)
	Login(user *LoginUser) (*Auth, error)
	Logout(token string) (string, error)
	GetOne(code string) (user *UserDto, errors error)
	GetAll(search string) ([]*UserDto, error)
	Update(code string, user *User) (*UserDto, error)
	PasswordUpdate(oldpassword, email, newpassword string) (string, string, error)
	Delete(code string) (string, error)
}
type userService struct {
	repo UserrepoInterface
}

func NewUserService(repository UserrepoInterface) UserServiceInterface {
	return &userService{
		repository,
	}
}
func (service *userService) Create(user *User) (*UserDto, error) {
	return service.repo.Create(user)
}

func (service *userService) GetAll(search string) ([]*UserDto, error) {
	return service.repo.GetAll(search)
}

func (service *userService) Login(auser *LoginUser) (*Auth, error) {
	return service.repo.Login(auser)
}

func (service *userService) Logout(token string) (string, error) {
	return service.repo.Logout(token)
}
func (service *userService) GetOne(code string) (*UserDto, error) {
	return service.repo.GetOne(code)
}
func (service *userService) PasswordUpdate(oldpassword, email, newpassword string) (string, string, error) {
	return service.repo.PasswordUpdate(oldpassword, email, newpassword)
}

func (service *userService) Update(code string, user *User) (*UserDto, error) {
	return service.repo.Update(code, user)
}

func (service *userService) Delete(id string) (string, error) {
	return service.repo.Delete(id)
}
