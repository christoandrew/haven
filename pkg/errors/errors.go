package errors

import "errors"

func InvalidDataError() error {
	return errors.New("invalid data")
}

func UserNotFoundError() error {
	return errors.New("user not found")
}

func InvalidCredentialsError() error {
	return errors.New("invalid credentials")
}
