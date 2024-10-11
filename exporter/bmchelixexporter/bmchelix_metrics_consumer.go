package bmchelixexporter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"
)

// BmcHelixMetricsConsumer is responsible for sending the metrics payload to the BMC Helix
type BmcHelixMetricsConsumer struct {
	authManager *AuthManager
	endpoint    string
	logger      *zap.Logger
	timeout     time.Duration
}

// SendHelixPayload sends the metrics payload to the BMC Helix
func (mc *BmcHelixMetricsConsumer) SendHelixPayload(ctx context.Context, payload []BmcHelixMetric) error {
	if (len(payload) == 0) {
		mc.logger.Warn("Payload is empty, nothing to send")
		return nil
	}


	// Get the authentication token using the AuthManager
	token, err := mc.authManager.GetToken(ctx)
	if err != nil {
		mc.logger.Error("Failed to get authentication token", zap.Error(err))
		return err
	}

	// Log the payload being sent
	mc.logger.Debug("Sending payload to BMC Helix payload", zap.Any("payload", payload))

	// Serialize the payload into JSON format
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		mc.logger.Error("Failed to marshal metrics payload", zap.Error(err))
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	url := mc.endpoint + "/metrics-gateway-service/api/v1.0/insert"

	// Create a new HTTP request to send the payload
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		mc.logger.Error("Failed to create HTTP request", zap.Error(err))
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	// Create an HTTP client with the specified timeout
	client := &http.Client{
		Timeout: mc.timeout,
	}

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		mc.logger.Error("Failed to send request to BMC Helix", zap.Error(err))
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		mc.logger.Error("Received non-2xx response from BMC Helix", zap.Int("status_code", resp.StatusCode))
		return fmt.Errorf("received non-2xx response: %d", resp.StatusCode)
	}

	mc.logger.Info("Successfully sent payload to BMC Helix", zap.String("endpoint", mc.endpoint))
	return nil
}
