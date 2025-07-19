package jobs

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/hibiken/asynq"
	"github.com/mnizarzr/dot-test/config"
	"github.com/mnizarzr/dot-test/utils"
)

// Job type constants
const (
	TypeEmailWelcome = "email:welcome"
)

// WelcomeEmailPayload represents the payload for welcome email job
type WelcomeEmailPayload struct {
	UserEmail string `json:"user_email"`
	UserName  string `json:"user_name"`
	UserRole  string `json:"user_role"`
	UserPassword string `json:"user_password"`
}

// EmailJobHandler handles email-related jobs
type EmailJobHandler struct {
	emailService *utils.EmailService
}

// NewEmailJobHandler creates a new email job handler
func NewEmailJobHandler(cfg *config.Config) *EmailJobHandler {
	return &EmailJobHandler{
		emailService: utils.NewEmailService(cfg),
	}
}

// NewWelcomeEmailTask creates a new welcome email task
func NewWelcomeEmailTask(userEmail, userName, userRole, userPlainPassword string) (*asynq.Task, error) {
	payload := WelcomeEmailPayload{
		UserEmail: userEmail,
		UserName:  userName,
		UserRole:  userRole,
		UserPassword: userPlainPassword,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(TypeEmailWelcome, payloadBytes), nil
}

// HandleWelcomeEmail processes welcome email job
func (h *EmailJobHandler) HandleWelcomeEmail(ctx context.Context, t *asynq.Task) error {
	var payload WelcomeEmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal welcome email payload: %w", err)
	}

	log.Printf("Sending welcome email to %s (%s)", payload.UserName, payload.UserEmail)

	err := h.emailService.SendWelcomeEmail(payload.UserEmail, payload.UserName, payload.UserRole, payload.UserPassword)
	if err != nil {
		return fmt.Errorf("failed to send welcome email to %s: %w", payload.UserEmail, err)
	}

	log.Printf("Welcome email sent successfully to %s", payload.UserEmail)
	return nil
}
