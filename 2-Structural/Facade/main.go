package main

import "fmt"

type IUser interface {
	Login(int, int) (*User, error)
	Register(int, int) (*User, error)
}

type User struct {
	Name string
}

func newUser(name string) IUser {
	return &User{
		Name: name,
	}
}

func (u *User) Login(phone, id int) (*User, error) {
	return &User{Name: "test login"}, nil
}

func (u *User) Register(phone, id int) (*User, error) {
	return &User{Name: "test register"}, nil
}

type UserService struct {
	User IUser
}

func NewUserService(name string) UserService {
	return UserService{
		User: newUser(name),
	}
}

func (us UserService) LoginOrRegister(phone, id int) (*User, error) {
	user, err := us.User.Login(phone, id)
	if err != nil {
		return nil, err
	}
	if user != nil {
		return user, nil
	}
	return us.User.Register(phone, id)
}

type IUserFacade interface {
	LoginOrRegister(int, int) (*User, error)
}

func main() {
	us := NewUserService("Gopher")
	res, _ := us.LoginOrRegister(123456789, 1)
	fmt.Println("result: ", res)
}
