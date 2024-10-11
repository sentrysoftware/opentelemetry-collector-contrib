package bmchelixexporter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// New instance of AuthManager for testing
func newTestAuthManager() *AuthManager {
	logger := zap.NewNop() // No-op logger for testing
	config := &Config{
		Endpoint: "https://example.com",
		Api: ApiConfig{
			AccessKey:       "test-access-key",
			AccessSecretKey: "test-secret-key",
			TenantId:        "test-tenant-id",
		},
	}
	return NewAuthManager(config, logger)
}

func TestAuthManager_GetToken_Success(t *testing.T) {
	t.Parallel()
	// Mock server to simulate BMC Helix API response
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authResponse := AuthResponse{
			JsonWebToken: "mock-jwt-token",
		}
		jsonResponse, _ := json.Marshal(authResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponse)
	}))
	defer mockServer.Close()

	// Initialize the AuthManager with the mock server's URL and a custom decodeJWTExpiry function
	authManager := newTestAuthManager()
	authManager.config.Endpoint = mockServer.URL
	authManager.decodeJWTExpiry = func(jwtToken string) (time.Time, error) {
		return time.Now().Add(1 * time.Hour), nil
	}

	// Test GetToken() method
	token, err := authManager.GetToken(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "mock-jwt-token", token)
	assert.True(t, time.Now().Before(authManager.expiry))
}

func TestAuthManager_GetToken_Failure(t *testing.T) {
	t.Parallel()
	// Mock server to simulate BMC Helix API failure
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer mockServer.Close()

	// Replace the real endpoint with the mock server's URL
	authManager := newTestAuthManager()
	authManager.config.Endpoint = mockServer.URL

	// Test GetToken() method when API returns a failure
	_, err := authManager.GetToken(context.Background())
	assert.Error(t, err)
	assert.Equal(t, "failed to authenticate with BMC Helix API", err.Error())
}

func TestDecodeJWTExpiry_Success(t *testing.T) {
	t.Parallel()

	expiryTime := time.Now().Add(1 * time.Hour).Unix()
	// Create a JWT token with a future expiry time
	jwtToken := createMockJWT(jwt.MapClaims{
		"exp": float64(expiryTime),
	})

	// Test decodeJWTExpiry() method
	expiry, err := decodeJWTExpiry(jwtToken)
	assert.NoError(t, err)
	assert.WithinDuration(t, time.Unix(expiryTime, 0), expiry, time.Second)
}

func TestDecodeJWTExpiry_Failure_On_No_Exp_Claim(t *testing.T) {
	t.Parallel()

	jwtToken := createMockJWT(jwt.MapClaims{
		"other": "claim",
	})

	// Test decodeJWTExpiry() method with invalid token
	expiry, err := decodeJWTExpiry(jwtToken)
	assert.Error(t, err)
	assert.True(t, expiry.IsZero())
}

// Helper function to create a mock JWT token with the given claim
func createMockJWT(claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte("secret"))
	return signedToken
}
