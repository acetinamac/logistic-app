package usecase

import (
	"errors"
	"logistics-app/backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Create(u *domain.User) error
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	DeleteByID(id uint) error
}

type UserService struct {
	repo UserRepo
}

func NewUserService(r UserRepo) *UserService {
	return &UserService{repo: r}
}

func (s *UserService) Register(email, password, fullName, phone string, role domain.Role) (*domain.User, error) {
	if email == "" || password == "" || fullName == "" {
		return nil, errors.New("email, password y full_name requeridos")
	}

	if role != domain.RoleClient && role != domain.RoleAdmin {
		role = domain.RoleClient
	}

	// Hash password before storing
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("no se pudo encriptar el password")
	}
	u := &domain.User{Email: email, Password: string(hash), FullName: fullName, Phone: phone, Role: role, IsActive: true}
	if err := s.repo.Create(u); err != nil {
		return nil, err
	}
	return u, nil
}

func (s *UserService) Delete(id uint) error {
	return s.repo.DeleteByID(id)
}

func (s *UserService) Authenticate(email, password string) (*domain.User, error) {
	if email == "" || password == "" {
		return nil, errors.New("email and password required")
	}

	u, err := s.repo.FindByEmail(email)
	if err != nil || u == nil {
		return nil, errors.New("invalid credentials")
	}

	// Compare provided password with stored bcrypt hash
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid password")
	}
	return u, nil
}

func (s *UserService) GetByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}
