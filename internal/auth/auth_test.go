package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	password1 := "Pa$$word"
	password2 := "2ManyTabs!"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)
	
	tests := []struct {
		name			string
		password		string
		hash			string
		wantErr			bool
		matchPassword	bool
	}{
		{
			name:			"Correct password",
			password:		password1,
			hash:			hash1,
			wantErr:		false,
			matchPassword:	true,
		},
		{
			name:			"Incorrect password",
			password:		"NotPassword1",
			hash:			hash1,
			wantErr:		false,
			matchPassword:	false,
		},
		{
			name:			"Different hash",
			password:		password1,
			hash:			hash2,
			wantErr:		false,
			matchPassword:	false,
		},
		{
			name:			"Empty password",
			password:		"",
			hash:			hash1,
			wantErr:		false,
			matchPassword:	false,
		},
		{
			name:			"Invalid hash",
			password:		password1,
			hash:			"hash1",
			wantErr:		true,
			matchPassword:	false,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s:\n\tCheckPasswordHash() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("%s:\n\tCheckPasswordHash() expects %v, got %v", tt.name, tt.matchPassword, match)
			}
		})
	}
}

func TestValidateJWT(t *testing.T) {
	testUserID := uuid.New()
	secretString := "SuperSecureString"
	shortExpiration, _ := time.ParseDuration("0s")
	longExpiration, _ := time.ParseDuration("5s")
	longTokenString, _ := MakeJWT(testUserID, secretString, longExpiration)
	shortTokenString, _ := MakeJWT(testUserID, secretString, shortExpiration)
	wrongTokenString, _ := MakeJWT(testUserID, "AnotherSecretString", longExpiration)
		
	tests := []struct {
		name			string
		userID			uuid.UUID
		secretString	string
		tokenString		string
		wantErr			bool
		wantValid		bool
	}{
		{
			name:			"Unexpired token",
			userID:			testUserID,
			secretString:	secretString,
			tokenString:	longTokenString,
			wantErr:		false,
			wantValid:		true,
		},
		{
			name:			"Expired token",
			userID:			testUserID,
			secretString:	secretString,
			tokenString:	shortTokenString,
			wantErr:		true,
			wantValid:		false,
		},
		{
			name:			"Wrong uuid",
			userID:			uuid.New(),
			secretString:	secretString,
			tokenString:	longTokenString,
			wantErr:		false,
			wantValid:		false,
		},
		{
			name:			"Wrong secret string",
			userID:			testUserID,
			secretString:	"NotTheSuperSecretString",
			tokenString:	longTokenString,
			wantErr:		true,
			wantValid:		false,
		},
		{
			name:			"Wrong token string",
			userID:			testUserID,
			secretString:	secretString,
			tokenString:	wrongTokenString,
			wantErr:		true,
			wantValid:		false,
		},
	}
		
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID, err := ValidateJWT(tt.tokenString, tt.secretString)
			if (err != nil) != tt.wantErr {
				t.Errorf("%s:\n\tValidateJWT() error = %v, wantErr %v", tt.name, err, tt.wantErr)
			}
			if (!tt.wantErr && (userID == tt.userID)) != tt.wantValid {
				t.Errorf("%s:\n\tValidateJWT() expects valid: %v with %v, got %v", tt.name, tt.wantValid, tt.userID, userID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
	want := "token123"
	req.Header.Set("Authorization", "Bearer " + want)
	got, err := GetBearerToken(req.Header)
	if err != nil {
		t.Errorf("GetBearerToken() got error: %v", err)
	}
	if got != want {
		t.Errorf("GetBearerToken() expected: %s, got: %s", want, got)
	}
}

func TestGetAPIKey(t *testing.T) {
	req, _ := http.NewRequest("GET", "https://api.example.com/data", nil)
	want := "apikey1234567890"
	req.Header.Set("Authorization", "ApiKey " + want)
	got, err := GetAPIKey(req.Header)
	if err != nil {
		t.Errorf("GetAPIKey() got error: %v", err)
	}
	if got != want {
		t.Errorf("GetAPIKey() expected: %s, got: %s", want, got)
	}
}