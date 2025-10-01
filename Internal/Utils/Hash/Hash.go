package HashPassword

import "golang.org/x/crypto/bcrypt"

func Hashpassword(password string) (hpassword string, error error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", err
	}
	return string(passwordBytes), nil
}

func Unhashpassword(plainpassword, hashedPassword string) (ispassmatch bool, erro error) {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainpassword))
	if err != nil {
		return false, err
	}
	return true, nil
}
