package interfaces

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/domain"
)

func TestLogin_Success(t *testing.T) {
	mockUserService := &application.MockUserService{
		AuthenticateFunc: func(email, password string) (*domain.User, error) {
			return &domain.User{
				ID:        1,
				Username:  "johndoe",
				Email:     "john@example.com",
				FirstName: "John",
				LastName:  "Doe",
				Role:      "user",
				Password:  "hashedpassword",
			}, nil
		},
	}
	handler := &HTTPHandler{UserService: mockUserService}

	reqBody := map[string]string{
		"email":    "john@example.com",
		"password": "password",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var respBody map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("response body unmarshal failed: %v", err)
	}

	if respBody["username"] != "johndoe" {
		t.Errorf("expected username %s, got %s", "johndoe", respBody["username"])
	}

	// Validate cookie
	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatal("expected a cookie, got none")
	}
	if cookies[0].Name != "token" {
		t.Errorf("expected cookie name %s, got %s", "token", cookies[0].Name)
	}
}

func TestRegister_Success(t *testing.T) {
	handler := &HTTPHandler{UserService: &application.MockUserService{
		RegisterUserFunc: func(user domain.CreateUserRequest) error {
			return nil
		},
	}}

	reqBody := domain.CreateUserRequest{
		Email:    "newuser@example.com",
		Password: "password123",
	}
	
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.RegisterUser(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}
}

func TestUpdateUser_Success(t *testing.T) {
	mockUser := &domain.User{
		ID:        1,
		Username:  "existinguser@example.com",
		Email:     "existinguser@example.com",
		FirstName: "Existing",
		LastName:  "User",
		Role:      "user",
	}
	handler := &HTTPHandler{UserService: &application.MockUserService{
		UpdateUserDataFunc: func(userID int64, email string, password string, firstName string, lastName string, bio string, profilePic string) error {
			return nil
		},
	}}

	handler.UserService = &application.MockUserService{
		AuthenticateFunc: func(email, password string) (*domain.User, error) {
			return mockUser, nil
		},
	}

	reqBody := map[string]string{
		"firstName": "Updated",
		"lastName":  "Name",
	}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest("PUT", "/update", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	handler.UpdateUser(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}
	fmt.Println(rec.Body.Bytes())
	var respBody map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("response body unmarshal failed: %v", err)
	}
	fmt.Println(respBody)
	// if respBody["firstName"] != "Updated" {
	// 	t.Errorf("expected firstName %s, got %s", "Updated", respBody["firstName"])
	// }
}

// TestGetUserProfile_Success tests retrieving the user profile with valid authorization
func TestGetUserProfile_Success(t *testing.T) {
	mockUser := &domain.User{
		ID:        1,
		Username:  "existinguser@example.com",
		Email:     "existinguser@example.com",
		FirstName: "Existing",
		LastName:  "User",
		Role:      "user",
	}
	// Arrange
	handler := &HTTPHandler{UserService: &application.MockUserService{
		GetUserByIDFunc: func(userID int64) (*domain.User, error) {
			return mockUser, nil
		},
	}}

	req := httptest.NewRequest("GET", "/profile", nil)
	req.Header.Set("Authorization", "Bearer valid-jwt-token")
	rec := httptest.NewRecorder()

	// Act
	handler.GetUserProfile(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var respBody map[string]interface{}
	if err := json.Unmarshal(rec.Body.Bytes(), &respBody); err != nil {
		t.Fatalf("response body unmarshal failed: %v", err)
	}

	if respBody["email"] != "user1@example.com" {
		t.Errorf("expected email %s, got %s", "user1@example.com", respBody["email"])
	}
}
