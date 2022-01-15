package auth

type AuthenticationError struct {
	reason string
}

func (e *AuthenticationError) Error() string {
	return "authentication failure: " + e.reason
}

var ErrPasswordDidNotMatch = AuthenticationError{"password did not match"}
var ErrTokenExpired = AuthenticationError{"token has expired"}
