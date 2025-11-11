package auth

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	testUserID := uuid.New()
	testTokenSecret := "anything"
	expiresDuration, _ := time.ParseDuration("20s")
	tokenString, err := MakeJWT(testUserID, testTokenSecret, expiresDuration)
	if err != nil || tokenString == "" {
		t.Errorf("Error creating JWT: %v", err)
	}
	userID, err := ValidateJWT(tokenString, testTokenSecret)
	if err != nil || tokenString == "" {
		t.Errorf("Error parsing JWT: %v", err)
	}
	if userID != testUserID {
		t.Errorf("User IDs do not match. Expected %s, got %s", testUserID, userID)
	}
}

func TestValidateJWTExpiration(t *testing.T) {
	testUserID := uuid.New()
	testTokenSecret := "anything"
	expiresDuration, _ := time.ParseDuration("1s")
	testWaitDuration, _ := time.ParseDuration("2s")
	tokenString, err := MakeJWT(testUserID, testTokenSecret, expiresDuration)
	if err != nil || tokenString == "" {
		t.Errorf("Error creating JWT: %v", err)
	}
	fmt.Println("Waiting for token to expire...")
	time.Sleep(testWaitDuration)
	fmt.Println("testing duration elapsed")

	_, err = ValidateJWT(tokenString, testTokenSecret)
	if err != nil {
		fmt.Printf("Error parsing JWT, token expired: %v", err)
	}
}

func TestValidateJWTInvalidToken(t *testing.T) {
	testUserID := uuid.New()
	testTokenSecret := "anything"
	expiresDuration, _ := time.ParseDuration("20s")
	tokenString, err := MakeJWT(testUserID, testTokenSecret, expiresDuration)
	wrongTokenString := "blah"
	if err != nil || tokenString == "" {
		t.Errorf("error creating jwt: %v", err)
	}
	_, err = ValidateJWT(wrongTokenString, testTokenSecret)
	if err == nil {
		t.Errorf("Error, wrong token string, should return an error")
	}
}

func TestValidateJWTInvalidSecret(t *testing.T) {
	testUserID := uuid.New()
	testTokenSecret := "anything"
	expiresDuration, _ := time.ParseDuration("20s")
	tokenString, err := MakeJWT(testUserID, testTokenSecret, expiresDuration)
	if err != nil || tokenString == "" {
		t.Errorf("error creating jwt: %v", err)
	}
	_, err = ValidateJWT(tokenString, "wrong secret")
	if err == nil {
		t.Errorf("wrong secret should return an error")
	}
}
