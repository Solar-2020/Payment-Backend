package middleware

type authClient interface {
	GetUserIDByCookie(sessionToken string) (userID int, err error)
	CompareSecret(inputSecret string) (err error)
}
