package util

import (
	"bufio"
	"fmt"
	"log"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/howeyc/gopass"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/SowinskiBraeden/DeviceBookingAPI/databases"
	"github.com/SowinskiBraeden/DeviceBookingAPI/models"
)

func ValidMailAddress(address string) (string, bool) {
	addr, err := mail.ParseAddress(address)
	if err != nil {
		return "", false
	}
	return addr.Address, true
}

func Confirm(s string) bool {
	r := bufio.NewReader(os.Stdin)

	fmt.Printf("%s [y/n]: ", s)
	res, err := r.ReadString('\n')
	if err != nil {
		log.Fatal(err)
	}

	return strings.ToLower(strings.TrimSpace(res))[0] == 'y'
}

func CreateDefaultAdmin(DB databases.UserDatabase) models.User {
	fmt.Println()
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("First Name: ")
	firstname, _ := reader.ReadString('\n')
	fmt.Print("Last Name: ")
	lastname, _ := reader.ReadString('\n')
	fmt.Print("Email: ")
	email, _ := reader.ReadString('\n')
	fmt.Print("Password: ")
	password, _ := gopass.GetPasswd()

	// Clear values of new lines and enter characters
	firstname = strings.ReplaceAll(firstname, "\n", "")
	lastname = strings.ReplaceAll(lastname, "\n", "")
	email = strings.ReplaceAll(email, "\n", "")
	firstname = strings.ReplaceAll(firstname, "\r", "")
	lastname = strings.ReplaceAll(lastname, "\r", "")
	email = strings.ReplaceAll(email, "\r", "")

	var admin models.User
	admin.Details.FirstName = firstname
	admin.Details.LastName = lastname
	admin.Details.Email = email

	pass := strings.TrimSuffix(string(password), "\n")
	admin.Details.Password = admin.HashPassword(pass)
	admin.Details.TempPassword = false

	var aid string
	for {
		aid = GenerateID(6)
		if ValidateID(aid, DB) {
			break
		}
	}
	admin.Details.UID = aid

	admin.Details.UserType = models.TypeSuperUser
	admin.Details.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	admin.Details.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	admin.ID = primitive.NewObjectID().Hex()

	return admin
}
