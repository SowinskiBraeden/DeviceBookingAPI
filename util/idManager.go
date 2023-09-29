package util

import (
	"context"
	"crypto/rand"
	"io"

	"github.com/SowinskiBraeden/DeviceBookingAPI/databases"
	"go.uber.org/zap"

	"go.mongodb.org/mongo-driver/bson"
)

var table = [...]byte{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}

func ValidateID(id string, DB databases.UserDatabase) bool { // true: valid id, false: id already in use

	dbResp, err := DB.Find(context.TODO(), bson.M{})
	if err != nil {
		zap.S().With(err).Error("failed to get users")
		return false
	}

	// If len == 0 then ID is not in use
	if len(dbResp) == 0 {
		return true
	}
	return false
}

func GenerateID(length int) string {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length {
		panic(err)
	}
	for i := 0; i < len(b); i++ {
		b[i] = table[int(b[i])%len(table)]
	}
	return string(b)
}
