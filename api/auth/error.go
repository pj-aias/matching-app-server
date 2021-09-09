package auth

type AuthenticationError struct {
	reason string
}

func (e *AuthenticationError) Error() string {
	return "authentication failure: " + e.reason
}
