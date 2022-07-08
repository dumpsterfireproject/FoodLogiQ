package service

type User struct {
	UserID  string
	Company string
}

type AuthenticationService interface {
	ValidateToken(string) *User
}

type AuthenticationServiceImpl struct {
}

func NewAuthenticationService() *AuthenticationServiceImpl {
	return &AuthenticationServiceImpl{}
}

func (s *AuthenticationServiceImpl) ValidateToken(token string) *User {
	return users[token]
}

var users = map[string]*User{
	"74edf612f393b4eb01fbc2c29dd96671": {"12345", "Acme"},
	"d88b4b1e77c70ba780b56032db1c259b": {"98765", "Ajax"},
}
