// Package echo provide implementations of custom functionality for the echo framework.
package echo

// Handler interface implementation
type Handler interface {
	Serve(e *Echo)
}
