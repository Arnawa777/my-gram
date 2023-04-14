package helpers

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (res string, err error) {
	salt := 10
	arrByte := []byte(password)
	hash, err := bcrypt.GenerateFromPassword(arrByte, salt)
	return string(hash), err
}

// get string argument h and p
// representing a hashed password and a plain text password
func PasswordValid(h, p string) bool {
	hash, pass := []byte(h), []byte(p)

	// It uses the bcrypt package to compare the hashed password with the plain text password, and returns a boolean indicating whether or not they match
	err := bcrypt.CompareHashAndPassword(hash, pass) // true / false
	//  If err is nil, it means that the password is valid
	// if
	return err == nil
}
