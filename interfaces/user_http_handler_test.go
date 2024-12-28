package interfaces

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/bandvov/social-media-go/application"
	"github.com/bandvov/social-media-go/domain"
	"github.com/bandvov/social-media-go/utils"
)

func TestLogin_Success(t *testing.T) {
	f := "John"
	l := "Doe"
	p := "pic"
	b := "bio"
	cn := "access_token"
	tn := time.Now()

	expectedUser := domain.User{
		ID:         1,
		Username:   nil,
		Email:      "john@example.com",
		FirstName:  &f,
		LastName:   &l,
		ProfilePic: &p,
		Bio:        &b,
		Role:       "user",
		PostsCount: 0,
		CreatedAt:  tn,
		UpdatedAt:  tn,
		Status:     "active",
	}

	// Marshal the user struct into JSON
	userJSON, err := json.Marshal(expectedUser)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	utils.JWTSecretKey = []byte("testsecret")

	tests := []struct {
		name            string
		requestBody     interface{}
		mockUserService application.UserServiceInterface
		expectedStatus  int
		expectedBody    string
		expectedCookie  *string
	}{
		{
			name: "Valid Login",
			requestBody: domain.CreateUserRequest{
				Email:    "john@example.com",
				Password: "password123",
			},
			mockUserService: &application.MockUserService{
				AuthenticateFunc: func(email, password string) (*domain.User, error) {
					return &expectedUser, nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   string(userJSON),
			expectedCookie: &cn,
		},
		{
			name:            "Invalid Request Body",
			requestBody:     `{"email": "test@example.com"`, // Malformed JSON
			mockUserService: &application.MockUserService{},
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "{\"message\": \"invalid request body\"}",
			expectedCookie:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up handler and request
			handler := UserHTTPHandler{
				UserService: tt.mockUserService,
			}

			var reqBody []byte
			if body, ok := tt.requestBody.(domain.CreateUserRequest); ok {
				reqBody, _ = json.Marshal(body)
			} else {
				reqBody = []byte(tt.requestBody.(string))
			}

			req := httptest.NewRequest(http.MethodPost, "/users/login", bytes.NewReader(reqBody))
			rec := httptest.NewRecorder()

			// Call the handler
			handler.Login(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Fatalf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			// Compare the response body
			if reflect.DeepEqual(rec.Body, []byte(tt.expectedBody)) {
				t.Errorf("expected body %s, got %s", tt.expectedBody, rec.Body.String())
			}

			// Validate cookie
			cookies := rec.Result().Cookies()
			if tt.expectedCookie != nil {
				if len(cookies) == 0 {
					t.Fatal("expected a cookie, got none")
				}
				if cookies[0].Name != *tt.expectedCookie {
					t.Errorf("expected cookie name %s, got %s", *tt.expectedCookie, cookies[0].Name)
				}
			} else if len(cookies) > 0 {
				t.Errorf("did not expect a cookie, but got %v", cookies[0].Name)
			}
		})
	}
}

func TestRegister(t *testing.T) {
	tests := []struct {
		name            string
		reqBody         domain.CreateUserRequest
		mockUserService application.UserServiceInterface
		expectedStatus  int
		expectedBody    string
	}{
		{
			name: "Success",
			reqBody: domain.CreateUserRequest{
				Email:    "newuser@example.com",
				Password: "password123",
			},
			mockUserService: &application.MockUserService{
				RegisterUserFunc: func(user domain.CreateUserRequest) error {
					return nil
				},
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `{"message": "User registered successfully"}`,
		},
		{
			name: "MissingEmail",
			reqBody: domain.CreateUserRequest{
				Password: "password123", // Email is missing
			},
			mockUserService: &application.MockUserService{
				RegisterUserFunc: func(user domain.CreateUserRequest) error {
					return nil
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Email is required"}`,
		},
		{
			name: "MissingPassword",
			reqBody: domain.CreateUserRequest{
				Email: "newuser@example.com", // Password is missing
			},
			mockUserService: &application.MockUserService{
				RegisterUserFunc: func(user domain.CreateUserRequest) error {
					return nil
				},
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error": "Password is required"}`,
		},
		{
			name: "InternalServerError",
			reqBody: domain.CreateUserRequest{
				Email:    "newuser@example.com",
				Password: "password123",
			},
			mockUserService: &application.MockUserService{
				RegisterUserFunc: func(user domain.CreateUserRequest) error {
					return errors.New(`{"message": "Internal server error"}`)
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message": "Internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Given: We set up the handler with the mock service
			handler := &UserHTTPHandler{
				UserService: &application.MockUserService{},
			}

			// When: We marshal the request body and create the HTTP request
			body, _ := json.Marshal(tt.reqBody)
			req := httptest.NewRequest("POST", "/users/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")

			// Set up a recorder to capture the HTTP response
			rec := httptest.NewRecorder()

			// Act: Call the handler's RegisterUser function to process the request
			handler.RegisterUser(rec, req)

			// Then: Assert that the response status code matches the expected status
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			// And: Assert that the response body matches the expected body
			bodyResponse := rec.Body.String()
			if bodyResponse != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, bodyResponse)
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		requestBody    map[string]string
		mockService    *application.MockUserService
		expectedStatus int
		expectedBody   string
		userId         int
	}{
		// {
		// 	name: "Successful Update",
		// 	requestBody: map[string]string{
		// 		"firstName": "Updated",
		// 		"lastName":  "Name",
		// 	},
		// 	mockService: &application.MockUserService{
		// 		UpdateUserDataFunc: func(user domain.User) error {
		// 			return nil
		// 		},
		// 	},
		// 	expectedStatus: http.StatusOK,
		// 	expectedBody:   `{"message":"User updated successfully"}`,
		// 	userId:         1,
		// },
		// {
		// 	name:        "Validation Error - Missing Fields",
		// 	requestBody: map[string]string{},
		// 	mockService: &application.MockUserService{
		// 		UpdateUserDataFunc: func(user domain.User) error {
		// 			return nil
		// 		},
		// 	},
		// 	expectedStatus: http.StatusBadRequest,
		// 	expectedBody:   `{"message":"First name and last name are required"}`,
		// 	userId:         1,
		// },
		// {
		// 	name: "Service Error",
		// 	requestBody: map[string]string{
		// 		"firstName": "Error",
		// 		"lastName":  "Case",
		// 	},
		// 	mockService: &application.MockUserService{
		// 		UpdateUserDataFunc: func(user domain.User) error {
		// 			return errors.New("failed to update user")
		// 		},
		// 	},
		// 	expectedStatus: http.StatusInternalServerError,
		// 	expectedBody:   `{"message":"Failed to update user"}`,
		// },
		{
			name: "missing user id",
			requestBody: map[string]string{
				"firstName": "Error",
				"lastName":  "Case",
			},
			mockService:    &application.MockUserService{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message": "invalid user ID"}`,
		},
	}

	// Iterate over test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Prepare request body
			body, _ := json.Marshal(tt.requestBody)

			// Create request and response recorder
			req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%v", tt.userId), bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			// Create handler with mock service
			handler := &UserHTTPHandler{UserService: tt.mockService}

			// Call the UpdateUser handler
			handler.UpdateUser(rec, req)

			// Assert status code
			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}

			// Assert response body
			if rec.Body.String() != tt.expectedBody {
				t.Errorf("expected response body '%s', got '%s'", tt.expectedBody, rec.Body.String())
			}
		})
	}
}

TestGetUserProfile_Success tests retrieving the user profile with valid authorization
func TestGetUserProfile_Success(t *testing.T) {
	utils.JWTSecretKey = []byte("test secret")

	f := "Existing"
	l := "User"
	u := "John"
	mockUser := &domain.User{
		ID:        1,
		Username:  &u,
		Email:     "existinguser@example.com",
		FirstName: &f,
		LastName:  &l,
		Role:      "user",
	}
	// Arrange
	handler := &UserHTTPHandler{UserService: &application.MockUserService{
		GetUserByIDFunc: func(userID int) (*domain.User, error) {
			return mockUser, nil
		},
	}}

	req := httptest.NewRequest("GET", "/users/1/profile", nil)
	req.Header.Set("Authorization", "Bearer valid-jwt-token")
	ctx := context.WithValue(req.Context(), userIDKey, 1)
	ctx = context.WithValue(ctx, isAdminKey, true)
	req = req.WithContext(ctx)
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
