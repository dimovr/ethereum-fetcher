package auth

import "ethereum_fetcher/internal/services/errors"

type AuthError struct {
	errors.ServiceError
}

func authError(msg string) AuthError {
	return AuthError{errors.NewServiceError(msg)}
}

type JwtError struct {
	AuthError
}

func jwtError(msg string) JwtError {
	return JwtError{authError(msg)}
}

var (
	TokenSigningFailed = jwtError("failed to sign token")
	InvalidSignature   = jwtError("invalid jwt signature")
	InvalidToken       = jwtError("invalid token")
	NoSubject          = jwtError("no subject found")
)

type UserError struct {
	AuthError
}

func userError(msg string) UserError {
	return UserError{authError(msg)}
}

var (
	UsernameNotFound = userError("username not found")
	InvalidPassword  = userError("invalid password")
)
