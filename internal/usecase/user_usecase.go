package usecase

import (
	"errors"

	"event-budaya-ticketing-bcc/config"
	"event-budaya-ticketing-bcc/internal/domain"
	"event-budaya-ticketing-bcc/internal/repository"
	"event-budaya-ticketing-bcc/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	Register(req *domain.CreateUserRequest) (*domain.UserResponse, error)
	Login(req *domain.LoginRequest) (*domain.LoginResponse, error)
	GetAllUsers() ([]domain.UserResponse, error)
	GetUserByID(id uint) (*domain.UserResponse, error)
	UpdateUser(id uint, req *domain.UpdateUserRequest) (*domain.UserResponse, error)
	DeleteUser(id uint) error
}

type userUsecase struct {
	userRepo repository.UserRepository
}

// NewUserUsecase creates a new instance of UserUsecase
func NewUserUsecase(userRepo repository.UserRepository) UserUsecase {
	return &userUsecase{userRepo: userRepo}
}

func (u *userUsecase) Register(req *domain.CreateUserRequest) (*domain.UserResponse, error) {
	// Check if email already exists
	existingUser, _ := u.userRepo.FindByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Set default role if not provided
	role := req.Role
	if role == "" {
		role = "user"
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Role:     role,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	response := user.ToResponse()
	return &response, nil
}

func (u *userUsecase) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := jwt.GenerateToken(user.ID, user.Email, user.Role, config.AppConfig.JWTSecret, config.AppConfig.JWTExpiryHours)
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	return &domain.LoginResponse{
		User:  user.ToResponse(),
		Token: token,
	}, nil
}

func (u *userUsecase) GetAllUsers() ([]domain.UserResponse, error) {
	users, err := u.userRepo.FindAll()
	if err != nil {
		return nil, errors.New("failed to fetch users")
	}

	var responses []domain.UserResponse
	for _, user := range users {
		responses = append(responses, user.ToResponse())
	}

	return responses, nil
}

func (u *userUsecase) GetUserByID(id uint) (*domain.UserResponse, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	response := user.ToResponse()
	return &response, nil
}

func (u *userUsecase) UpdateUser(id uint, req *domain.UpdateUserRequest) (*domain.UserResponse, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// Update fields if provided
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Email != "" {
		// Check if new email already exists (and not the same user)
		existingUser, _ := u.userRepo.FindByEmail(req.Email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already in use")
		}
		user.Email = req.Email
	}
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, errors.New("failed to hash password")
		}
		user.Password = string(hashedPassword)
	}
	if req.Role != "" {
		user.Role = req.Role
	}

	if err := u.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update user")
	}

	response := user.ToResponse()
	return &response, nil
}

func (u *userUsecase) DeleteUser(id uint) error {
	_, err := u.userRepo.FindByID(id)
	if err != nil {
		return errors.New("user not found")
	}

	if err := u.userRepo.Delete(id); err != nil {
		return errors.New("failed to delete user")
	}

	return nil
}
