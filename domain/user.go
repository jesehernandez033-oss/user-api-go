package domain

import (
	"net/mail"
	"time"
	"unicode"
)

type User struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Birthday    string `json:"birthday"`
	Address     string `json:"address"`
}

// 📧 Validar email
func (u *User) IsValidEmail() bool {
	_, err := mail.ParseAddress(u.Email)
	return err == nil
}

// 🔐 Validar password
func (u *User) IsValidPassword() bool {
	if len(u.Password) < 8 || len(u.Password) > 16 {
		return false
	}

	var hasUpper, hasNumber, hasSpecial bool

	for _, char := range u.Password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsDigit(char):
			hasNumber = true
		case !unicode.IsLetter(char) && !unicode.IsDigit(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasNumber && hasSpecial
}

// 📅 Validar edad
func (u *User) IsAdult() bool {
	birthDate, err := time.Parse("2006-01-02", u.Birthday)
	if err != nil {
		return false
	}

	now := time.Now()
	age := now.Year() - birthDate.Year()
	if now.YearDay() < birthDate.YearDay() {
		age--
	}

	return age >= 14
}
