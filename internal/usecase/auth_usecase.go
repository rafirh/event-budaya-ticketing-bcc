package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"event-budaya-ticketing-bcc/internal/domain"
	"event-budaya-ticketing-bcc/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(req *domain.RegisterRequest) (*domain.UserResponse, error)
	Login(req *domain.LoginRequest) (*domain.LoginResponse, error)
	GetMe(id string) (*domain.UserResponse, error)
	Logout(token string) error
}

type authUsecase struct {
	userRepo  repository.UserRepository
	tokenRepo repository.PersonalAccessTokenRepository
}

func NewAuthUsecase(userRepo repository.UserRepository, tokenRepo repository.PersonalAccessTokenRepository) AuthUsecase {
	return &authUsecase{userRepo: userRepo, tokenRepo: tokenRepo}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (u *authUsecase) Register(req *domain.RegisterRequest) (*domain.UserResponse, error) {
	existing, _ := u.userRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	role := req.Role
	if role == "" {
		role = "user"
	}

	user := &domain.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Role:     role,
		Gender:   req.Gender,
	}

	if err := u.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (u *authUsecase) Login(req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	rawToken, err := generateToken()
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	pat := &domain.PersonalAccessToken{
		UserID: user.ID,
		Token:  rawToken,
		Name:   "default",
	}
	if err := u.tokenRepo.Create(pat); err != nil {
		return nil, errors.New("failed to save token")
	}

	return &domain.LoginResponse{
		User:  user.ToResponse(),
		Token: rawToken,
	}, nil
}

func (u *authUsecase) GetMe(id string) (*domain.UserResponse, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	resp := user.ToResponse()
	return &resp, nil
}

func (u *authUsecase) Logout(token string) error {
	return u.tokenRepo.DeleteByToken(token)
}
