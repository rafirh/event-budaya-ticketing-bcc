package usecase

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"mime/multipart"
	"net/url"
	"strings"
	"time"

	"event-budaya-ticketing-bcc/internal/dto"
	"event-budaya-ticketing-bcc/internal/model"
	"event-budaya-ticketing-bcc/internal/repository"
	"event-budaya-ticketing-bcc/pkg/email"
	"event-budaya-ticketing-bcc/pkg/oauth"
	"event-budaya-ticketing-bcc/pkg/storage"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthUsecase interface {
	Register(req *dto.RegisterRequest) (*dto.UserResponse, error)
	Login(req *dto.LoginRequest) (*dto.LoginResponse, error)
	VerifyEmail(token string) error
	ResendVerificationEmail(email string) error
	GetMe(id string) (*dto.UserResponse, error)
	UpdateProfile(userID string, req *dto.UpdateProfileRequest, photo *multipart.FileHeader) (*dto.UserResponse, error)
	ChangePassword(userID string, req *dto.ChangePasswordRequest) error
	Logout(token string) error
	GoogleLoginURL(state string) string
	GoogleCallback(ctx context.Context, code, state string) (*dto.LoginResponse, error)
}

type authUsecase struct {
	userRepo              repository.UserRepository
	tokenRepo             repository.PersonalAccessTokenRepository
	emailVerificationRepo repository.EmailVerificationTokenRepository
	mailSender            email.Sender
	appURL                string
	uploader              storage.Uploader
	googleOAuthProvider   *oauth.GoogleOAuthProvider
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	tokenRepo repository.PersonalAccessTokenRepository,
	emailVerificationRepo repository.EmailVerificationTokenRepository,
	mailSender email.Sender,
	appURL string,
	uploader storage.Uploader,
	googleOAuthProvider *oauth.GoogleOAuthProvider,
) AuthUsecase {
	return &authUsecase{
		userRepo:              userRepo,
		tokenRepo:             tokenRepo,
		emailVerificationRepo: emailVerificationRepo,
		mailSender:            mailSender,
		appURL:                appURL,
		uploader:              uploader,
		googleOAuthProvider:   googleOAuthProvider,
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashToken(rawToken string) string {
	h := sha256.Sum256([]byte(rawToken))
	return hex.EncodeToString(h[:])
}

func (u *authUsecase) buildVerificationLink(rawToken string) string {
	base := strings.TrimRight(strings.TrimSpace(u.appURL), "/")
	if base == "" {
		base = "http://localhost:3000"
	}

	return fmt.Sprintf("%s/api/auth/verify-email?token=%s", base, url.QueryEscape(rawToken))
}

func (u *authUsecase) sendVerificationEmail(toName, toEmail, verificationLink string) error {
	if u.mailSender == nil {
		return errors.New("email sender is not configured")
	}

	subject := "Aktivasi Akun LokaBudaya"
	body := fmt.Sprintf(
		"<p>Halo %s,</p><p>Terima kasih sudah mendaftar di LokaBudaya.</p><p>Silakan aktivasi akun Anda dengan klik link berikut:</p><p><a href=\"%s\">Aktivasi Akun</a></p><p>Link ini hanya berlaku 15 menit.</p>",
		toName,
		verificationLink,
	)

	if err := u.mailSender.Send(toEmail, subject, body); err != nil {
		return errors.New("failed to send verification email: " + err.Error())
	}

	return nil
}

func (u *authUsecase) createAndSendVerificationToken(user *model.User) error {
	rawToken, err := generateToken()
	if err != nil {
		return errors.New("failed to create verification token")
	}

	if err := u.emailVerificationRepo.DeleteByUserID(user.ID); err != nil {
		return errors.New("failed to create verification token")
	}

	verificationToken := &model.EmailVerificationToken{
		UserID:    user.ID,
		TokenHash: hashToken(rawToken),
		ExpiresAt: time.Now().Add(15 * time.Minute),
	}

	if err := u.emailVerificationRepo.Create(verificationToken); err != nil {
		return errors.New("failed to create verification token")
	}

	verificationLink := u.buildVerificationLink(rawToken)
	if err := u.sendVerificationEmail(user.Name, user.Email, verificationLink); err != nil {
		_ = u.emailVerificationRepo.DeleteByUserID(user.ID)
		return err
	}

	return nil
}

func (u *authUsecase) Register(req *dto.RegisterRequest) (*dto.UserResponse, error) {
	existing, err := u.userRepo.FindByEmail(req.Email)
	if err == nil && existing != nil {
		if existing.EmailVerifiedAt != nil {
			return nil, errors.New("email already registered")
		}

		if err := u.tokenRepo.DeleteByUserID(existing.ID); err != nil {
			return nil, errors.New("failed to refresh unverified account")
		}
		if err := u.emailVerificationRepo.DeleteByUserID(existing.ID); err != nil {
			return nil, errors.New("failed to refresh unverified account")
		}
		if err := u.userRepo.DeleteByID(existing.ID.String()); err != nil {
			return nil, errors.New("failed to refresh unverified account")
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("failed to check existing email")
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

	if err := u.createAndSendVerificationToken(user); err != nil {
		_ = u.userRepo.DeleteByID(user.ID.String())
		return nil, err
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

	if user.EmailVerifiedAt == nil {
		return nil, errors.New("email is not verified, please activate your account from the email link")
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

func (u *authUsecase) VerifyEmail(token string) error {
	token = strings.TrimSpace(token)
	if token == "" {
		return errors.New("verification token is required")
	}

	tokenHash := hashToken(token)
	verificationToken, err := u.emailVerificationRepo.FindByTokenHash(tokenHash)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("invalid or expired verification token")
		}
		return errors.New("failed to verify email")
	}

	if verificationToken.UsedAt != nil {
		return errors.New("verification link has already been used")
	}

	if time.Now().After(verificationToken.ExpiresAt) {
		return errors.New("verification link has expired")
	}

	user, err := u.userRepo.FindByID(verificationToken.UserID.String())
	if err != nil {
		return errors.New("user not found")
	}

	if user.EmailVerifiedAt != nil {
		_ = u.emailVerificationRepo.DeleteByUserID(user.ID)
		return nil
	}

	now := time.Now()
	user.EmailVerifiedAt = &now
	if err := u.userRepo.Update(user); err != nil {
		return errors.New("failed to activate account")
	}

	if err := u.emailVerificationRepo.DeleteByUserID(user.ID); err != nil {
		return errors.New("failed to finalize email verification")
	}

	return nil
}

func (u *authUsecase) ResendVerificationEmail(email string) error {
	email = strings.TrimSpace(email)
	if email == "" {
		return errors.New("email is required")
	}

	user, err := u.userRepo.FindByEmail(email)
	if err != nil {
		return errors.New("user not found")
	}

	if user.EmailVerifiedAt != nil {
		return errors.New("account is already verified")
	}

	latestToken, err := u.emailVerificationRepo.FindLatestByUserID(user.ID)
	if err == nil {
		if latestToken.UsedAt == nil && time.Now().Before(latestToken.ExpiresAt) {
			return errors.New("verification link is still active, please use the existing link")
		}
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("failed to resend verification email")
	}

	if err := u.createAndSendVerificationToken(user); err != nil {
		return err
	}

	return nil
}

func (u *authUsecase) GetMe(id string) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("user not found")
	}

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (u *authUsecase) UpdateProfile(userID string, req *dto.UpdateProfileRequest, photo *multipart.FileHeader) (*dto.UserResponse, error) {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.Name != nil {
		user.Name = *req.Name
	}
	if req.Phone != nil {
		user.Phone = req.Phone
	}
	if req.Gender != nil {
		user.Gender = req.Gender
	}

	if photo != nil {
		if u.uploader == nil {
			return nil, errors.New("storage uploader is not configured")
		}

		url, uploadErr := u.uploader.UploadImage(context.Background(), photo, "profiles")
		if uploadErr != nil {
			return nil, errors.New("failed to upload profile photo: " + uploadErr.Error())
		}
		user.ProfilePhoto = &url
	}

	if err := u.userRepo.Update(user); err != nil {
		return nil, errors.New("failed to update profile")
	}

	resp := dto.ToUserResponse(user)
	return &resp, nil
}

func (u *authUsecase) ChangePassword(userID string, req *dto.ChangePasswordRequest) error {
	user, err := u.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return errors.New("current password is incorrect")
	}

	if req.CurrentPassword == req.NewPassword {
		return errors.New("new password must be different from current password")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user.Password = string(hashedPassword)
	if err := u.userRepo.Update(user); err != nil {
		return errors.New("failed to update password")
	}

	return nil
}

func (u *authUsecase) Logout(token string) error {
	return u.tokenRepo.DeleteByToken(token)
}

func (u *authUsecase) GoogleLoginURL(state string) string {
	return u.googleOAuthProvider.GetAuthURL(state)
}

func (u *authUsecase) GoogleCallback(ctx context.Context, code, state string) (*dto.LoginResponse, error) {
	googleUser, err := u.googleOAuthProvider.GetUserInfo(ctx, code)
	if err != nil {
		return nil, err
	}

	if googleUser.Email == "" {
		return nil, errors.New("email not provided by Google")
	}

	user, err := u.userRepo.FindByEmail(googleUser.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newUser := &model.User{
				Email:           googleUser.Email,
				Name:            googleUser.Name,
				Password:        "", // Will be set below
				Role:            "user",
				EmailVerifiedAt: &time.Time{},
				CreatedAt:       time.Now(),
			}

			if randPass, err := generateRandomPassword(); err == nil {
				hashPassword, _ := bcrypt.GenerateFromPassword([]byte(randPass), bcrypt.DefaultCost)
				newUser.Password = string(hashPassword)
			}

			now := time.Now()
			newUser.EmailVerifiedAt = &now

			if err := u.userRepo.Create(newUser); err != nil {
				return nil, errors.New("failed to create user from Google account")
			}

			user = newUser
		} else {
			return nil, errors.New("failed to get user")
		}
	} else {
		if user.EmailVerifiedAt == nil {
			now := time.Now()
			user.EmailVerifiedAt = &now
			if err := u.userRepo.Update(user); err != nil {
				return nil, errors.New("failed to update user email verification status")
			}
		}
	}

	rawToken, err := generateToken()
	if err != nil {
		return nil, errors.New("failed to generate token")
	}

	pat := &model.PersonalAccessToken{
		UserID: user.ID,
		Token:  rawToken,
		Name:   "google_login",
	}
	if err := u.tokenRepo.Create(pat); err != nil {
		return nil, errors.New("failed to save token")
	}

	return &dto.LoginResponse{
		User:  dto.ToUserResponse(user),
		Token: rawToken,
	}, nil
}

func generateRandomPassword() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
