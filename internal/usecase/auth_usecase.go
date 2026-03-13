package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"errors"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase interface {
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	GetMe(id string) (*dto.UserResponse, error)
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

func (u *authUsecase) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
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

	user := &model.User{
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

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (u *authUsecase) Login(req *dto.LoginRequest) (*dto.LoginResponse, error) {
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

	pat := &model.PersonalAccessToken{
		UserID: user.ID,
		Token:  rawToken,
		Name:   "default",
	}
	if err := u.tokenRepo.Create(pat); err != nil {
		return nil, errors.New("failed to save token")
	}

	return &dto.LoginResponse{
		User:  dto.ToUserResponse(user),
		Token: rawToken,
	}, nil
}

func (u *authUsecase) GetMe(id string) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (u *authUsecase) Logout(token string) error {
	return u.tokenRepo.DeleteByToken(token)
}
