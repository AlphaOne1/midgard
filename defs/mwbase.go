// SPDX-FileCopyrightText: 2026 The midgard contributors.
// SPDX-License-Identifier: MPL-2.0

// Package defs contains the common types and functions for all midgard handlers.
package defs

import (
	"errors"
	"log/slog"
	"net/http"
	"reflect"
)

// ErrNilLogger is returned when a logger is nil.
var ErrNilLogger = errors.New("logger cannot be nil")

// ErrNilHandler is returned when a handler is nil.
var ErrNilHandler = errors.New("handler cannot be nil")

// ErrNotInitialized is returned when a middleware is used before it has been initialized.
var ErrNotInitialized = errors.New("middleware not initialized")

// MWBase contains the basic middleware information common to each midgard handler.
type MWBase struct {
	log      *slog.Logger // logger
	logLevel slog.Level   // logLevel
	next     http.Handler // next contains the next handler in the handler chain.
}

// MWBaser is the interface used to get the basic middleware information as defined in MWBase.
type MWBaser interface {
	// GetMWBase gets the MWBase structure out of a midgard handler.
	GetMWBase() *MWBase
}

// Log gets the configured slog.Logger.
func (mw *MWBase) Log() *slog.Logger {
	if mw != nil && mw.log != nil {
		return mw.log
	}

	return slog.Default()
}

// SetLog sets a new slog.Logger to use for logging.
func (mw *MWBase) SetLog(l *slog.Logger) error {
	if mw == nil {
		return ErrNotInitialized
	}

	if l == nil {
		return ErrNilLogger
	}

	mw.log = l

	return nil
}

// LogLevel gets the currently configured log level.
func (mw *MWBase) LogLevel() slog.Level {
	if mw != nil {
		return mw.logLevel
	}

	return slog.LevelInfo
}

// SetLogLevel sets the new log level to use.
func (mw *MWBase) SetLogLevel(l slog.Level) error {
	if mw != nil {
		mw.logLevel = l

		return nil
	}

	return ErrNotInitialized
}

// Next gets the next handler in a chain of handlers.
func (mw *MWBase) Next() http.Handler {
	if mw != nil {
		return mw.next
	}

	return nil
}

// SetNext sets the next handler for the chain of handlers.
func (mw *MWBase) SetNext(n http.Handler) error {
	if mw == nil {
		return ErrNotInitialized
	}

	if n == nil {
		return ErrNilHandler
	}

	mw.next = n

	return nil
}

// WithLogger is a convenience function to easily write the functional options
// for each handler.
func WithLogger[T MWBaser](l *slog.Logger) func(h T) error {
	return func(h T) error {
		if value := reflect.ValueOf(h); !value.IsValid() || value.IsNil() {
			return ErrNilHandler
		}

		return h.GetMWBase().SetLog(l)
	}
}

// WithLogLevel is a convenience function to easily write the functional options
// for each handler.
func WithLogLevel[T MWBaser](l slog.Level) func(h T) error {
	return func(h T) error {
		if value := reflect.ValueOf(h); !value.IsValid() || value.IsNil() {
			return ErrNilHandler
		}

		return h.GetMWBase().SetLogLevel(l)
	}
}
