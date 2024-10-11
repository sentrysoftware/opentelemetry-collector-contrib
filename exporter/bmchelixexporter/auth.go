package bmchelixexporter

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

type AuthResponse struct {
	UserId       string `json:"user_id"`
	PrincipalId  string `json:"principal_id"`
	TenantId     string `json:"tenant_id"`
	TenantName   string `json:"tenant_name"`
	Token        string `json:"token"`
	JsonWebToken string `json:"json_web_token"`
}

type AuthManager struct {
	config          *Config
	logger          *zap.Logger
	token           string
	tokenLock       sync.RWMutex
	expiry          time.Time
	decodeJWTExpiry func(jwtToken string) (time.Time, error)
}

// NewAuthManager creates a new instance of the AuthManager
func NewAuthManager(config *Config, logger *zap.Logger) *AuthManager {
	return &AuthManager{
		config:          config,
		logger:          logger,
		decodeJWTExpiry: decodeJWTExpiry,
	}
}

// GetToken fetches the token if expired or not present, otherwise returns the stored token
func (a *AuthManager) GetToken(ctx context.Context) (string, error) {
	a.tokenLock.RLock()
	if a.token != "" && time.Now().Before(a.expiry) {
		a.tokenLock.RUnlock()
		return a.token, nil
	}
	a.tokenLock.RUnlock()

	// Fetch new token
	a.tokenLock.Lock()
	defer a.tokenLock.Unlock()

	if err := a.fetchToken(ctx); err != nil {
		return "", err
	}

	return a.token, nil
}

// fetchToken makes the REST API call to get a new JWT token
func (auth *AuthManager) fetchToken(ctx context.Context) error {
	url := auth.config.Endpoint + "/ims/api/v1/access_keys/login"

	body := map[string]string{
		"access_key":        auth.config.Api.AccessKey,
		"access_secret_key": auth.config.Api.AccessSecretKey,
		"tenant_id":         auth.config.Api.TenantId,
	}

	jsonBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("failed to authenticate with BMC Helix API")
	}

	var authResponse AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		return err
	}

	auth.token = authResponse.JsonWebToken
	expiry, err := auth.decodeJWTExpiry(authResponse.JsonWebToken) // assuming you decode the JWT to get expiry time

	if err != nil {
		auth.logger.Error("Failed to decode JWT", zap.Error(err))
		return err
	}

	auth.expiry = expiry

	auth.logger.Debug("Successfully authenticated with BMC Helix API", zap.String("token", auth.token), zap.Time("expiry", auth.expiry))

	return nil
}

// decodeJWTExpiry decodes the JWT and extracts the expiry time
func decodeJWTExpiry(jwtToken string) (time.Time, error) {
	token, _ := jwt.Parse(jwtToken, nil)

	if claims, ok := token.Claims.(jwt.MapClaims); ok && claims["exp"] != nil {
		exp := int64(claims["exp"].(float64))
		return time.Unix(exp, 0), nil
	}
	return time.Time{}, fmt.Errorf("failed to decode JWT")
}
