package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type GoogleUser struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Picture   string `json:"picture"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
}

type GoogleOAuthProvider struct {
	config *oauth2.Config
}

func NewGoogleOAuthProvider(clientID, clientSecret, redirectURL string) *GoogleOAuthProvider {
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	return &GoogleOAuthProvider{
		config: config,
	}
}

func (g *GoogleOAuthProvider) GetAuthURL(state string) string {
	return g.config.AuthCodeURL(state)
}

func (g *GoogleOAuthProvider) GetUserInfo(ctx context.Context, code string) (*GoogleUser, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return nil, errors.New("failed to exchange authorization code: " + err.Error())
	}

	client := g.config.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, errors.New("failed to get user info: " + err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to get user info: HTTP " + resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("failed to read response: " + err.Error())
	}

	var userInfo GoogleUser
	if err := json.Unmarshal(body, &userInfo); err != nil {
		return nil, errors.New("failed to parse user info: " + err.Error())
	}

	return &userInfo, nil
}
