package httpserver

import (
	"testing"

	"github.com/ElvisReis2K/Form-Builder/backend/internal/auth"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/config"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/forms"
	"github.com/ElvisReis2K/Form-Builder/backend/internal/responses"
)

func TestNewRegistersRoutes(t *testing.T) {
	server := New(
		config.Config{Address: "localhost:0"},
		nil,
		auth.NewHandler(nil, false, "", nil),
		forms.NewHandler(nil, nil),
		responses.NewHandler(nil, nil),
	)
	if server.Handler == nil {
		t.Fatal("expected server handler")
	}
}
