package healthcheck

import (
	"context"
	"net/http"

	"github.com/heptiolabs/healthcheck"
)

type health struct {
	healthcheck.Handler
	server *http.Server
}

// healthCheckable constrains the type accepted to the registerHealthChecks to
// instances that implement a Readiness and a Liveness probe
type healthCheckable interface {
	ReadinessCheck() (string, func() error)
	LivenessCheck() (string, func() error)
}

// New initiates the healthchecks as an HTTP server
func New(healthCheckables ...healthCheckable) *health { // nolint
	h := &health{Handler: healthcheck.NewHandler()}
	h.Append(healthCheckables...)
	return h
}

// Append will register the healthchecks
func (h *health) Append(healthCheckables ...healthCheckable) {
	for _, check := range healthCheckables {
		key, f := check.LivenessCheck()
		h.AddLivenessCheck(key, f)

		key, f = check.ReadinessCheck()
		h.AddReadinessCheck(key, f)
	}
}

// Start binds to the given address or returns an error
// Will block so start in a go routine.
func (h *health) Start(address string) error {
	h.server = &http.Server{Addr: address, Handler: h.Handler}

	if err := h.server.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}

	return nil
}

// Shutdown will close the underlying http server
func (h *health) Shutdown(ctx context.Context) error {
	return h.server.Shutdown(ctx)
}
