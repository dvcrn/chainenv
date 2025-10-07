package backend

import (
	"github.com/dvcrn/chainenv/logger"
)

// Backend defines the interface for password storage backends
type Backend interface {
	GetPassword(account string) (string, error)
	SetPassword(account, password string, update bool) error
	List() ([]string, error)
	GetMultiplePasswords(accounts []string) (map[string]string, error)
}

// BackendOpts contains options for configuring a backend
type BackendOpts struct {
	logger *logger.Logger
}

// BackendOption defines a function that can modify BackendOpts
type BackendOption func(*BackendOpts)

func WithLogger(logger *logger.Logger) BackendOption {
	return func(opts *BackendOpts) {
		opts.logger = logger
	}
}

// newBackendOpts creates a new BackendOpts with the given options applied
func newBackendOpts(options ...BackendOption) *BackendOpts {
	opts := &BackendOpts{
		logger: logger.NewLogger(false),
	}

	for _, opt := range options {
		opt(opts)
	}
	return opts
}
