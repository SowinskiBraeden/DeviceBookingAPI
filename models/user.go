package models

import (
	"math/rand"
	"strings"
	"time"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

const (
	lowerCharSet   = "abcdedfghijklmnopqrst"
	upperCharSet   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	specialCharSet = "!@#$%&*?,.`~"
	numberSet      = "0123456789"
	allCharSet     = lowerCharSet + upperCharSet + specialCharSet + numberSet
)

// User types
const (
	// SuperUsers can:
	// - Create/Delete/Modify businesses
	// - Create/Delete/Modify Admins & Users
	// - Assign Admins & Users to businesses
	TypeSuperUser int = 1

	// Admins can:
	// - Manage/Modify assigned business
	// - Use reservation system for their business
	// - Manage/Modify Cows & Devices for their business
	// - Create Users & assign to their business
	// - Promote Users to Admins in their business
	TypeAdmin = 2

	// Users can:
	// - Use reservation system for their business
	TypeUser = 3
)

type User struct {
	ID      string      `bson:"_id"`
	Details UserDetails `json:"details"`
}

type UserDetails struct {
	FirstName    string    `json:"firstname" validate:"required"`
	LastName     string    `json:"lastname" validate:"required"`
	Email        string    `json:"email" validate:"required"`
	Password     string    `json:"-" validate:"min=10,max=32"`
	TempPassword bool      `json:"temppassword"`
	UID          string    `json:"uid"`
	Business     string    `json:"business"`
	UserType     int       `json:"usertype"`
	Created_at   time.Time `json:"created_at"`
	Updated_at   time.Time `json:"updated_at"`
}

func (s *User) HashPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(hash)
}

func (a *User) ComparePasswords(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(a.Details.Password), []byte(password))
	return err == nil
}

func (a *User) CheckPasswordStrength(password string) bool {

	var hasUpper bool = false
	for _, r := range password {
		if unicode.IsUpper(r) && unicode.IsLetter(r) {
			hasUpper = true
		}
	}

	var hasLower bool = false
	for _, r := range password {
		if !unicode.IsLower(r) && unicode.IsLetter(r) {
			hasLower = true
		}
	}

	if strings.ContainsAny(password, specialCharSet) && hasLower && hasUpper && len(password) >= 8 {
		return true
	} else {
		return false
	}
}

func (a *User) GeneratePassword(passwordLength, minSpecialChar, minNum, minUpperCase int) string {
	var password strings.Builder

	//Set special character
	for i := 0; i < minSpecialChar; i++ {
		random := rand.Intn(len(specialCharSet))
		password.WriteString(string(specialCharSet[random]))
	}

	//Set numeric
	for i := 0; i < minNum; i++ {
		random := rand.Intn(len(numberSet))
		password.WriteString(string(numberSet[random]))
	}

	//Set uppercase
	for i := 0; i < minUpperCase; i++ {
		random := rand.Intn(len(upperCharSet))
		password.WriteString(string(upperCharSet[random]))
	}

	remainingLength := passwordLength - minSpecialChar - minNum - minUpperCase
	for i := 0; i < remainingLength; i++ {
		random := rand.Intn(len(allCharSet))
		password.WriteString(string(allCharSet[random]))
	}
	inRune := []rune(password.String())
	rand.Shuffle(len(inRune), func(i, j int) {
		inRune[i], inRune[j] = inRune[j], inRune[i]
	})
	return string(inRune)
}
