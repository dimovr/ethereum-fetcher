package handlers

import (
	"net/http"

	"ethereum_fetcher/api"
	"ethereum_fetcher/internal/services/auth"
	txnerrors "ethereum_fetcher/internal/services/transactions/types"
)

func mapError(err error) api.Error {
	return api.Error{Msg: err.Error()}
}

func toStatusCode(err error) int {
	// Authentication Errors
	if err == auth.InvalidToken || err == auth.NoSubject {
		return http.StatusUnauthorized
	}
	if err == auth.InvalidSignature {
		return http.StatusForbidden
	}
	if err == auth.TokenSigningFailed {
		return http.StatusInternalServerError
	}

	// User Errors
	if err == auth.UsernameNotFound || err == auth.InvalidPassword {
		return http.StatusBadRequest
	}

	// Transaction Errors
	if err == txnerrors.FailedToFetchTransaction {
		return http.StatusNotFound
	}

	// RLP Decoding Errors
	if err == txnerrors.InvalidHexEncoding ||
		err == txnerrors.InvalidRlpEncoding ||
		err == txnerrors.InvalidHashLength {
		return http.StatusBadRequest
	}

	// Default error handling
	return http.StatusInternalServerError
}
