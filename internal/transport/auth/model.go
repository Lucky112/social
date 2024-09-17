package auth

// Структура HTTP-запроса на регистрацию пользователя
type registerRequest struct {
	Email    string `json:"email"    validate:"required"`
	Login    string `json:"login"    validate:"required"`
	Password string `json:"password" validate:"required"`
}
type registerResponse struct {
	Id string `json:"id"`
}
type registerError struct {
	Message string `json:"msg"`
}

// Структура HTTP-запроса на вход в аккаунт
type loginRequest struct {
	Login    string `json:"login"    validate:"required"`
	Password string `json:"password" validate:"required"`
}

// Структура HTTP-ответа на вход в аккаунт
// В ответе содержится JWT-токен авторизованного пользователя
type loginResponse struct {
	AccessToken string `json:"access_token"`
}

type loginError struct {
	Message string `json:"msg"`
}
