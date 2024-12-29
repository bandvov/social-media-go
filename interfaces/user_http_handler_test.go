package interfaces

import (
	"bytes"
	"context"
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
	"github.com/lib/pq"
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

func TestRegisterUser(t *testing.T) {
	tests := []struct {
		name               string
		inputBody          interface{}
		mockResponse       error
		expectedStatusCode int
		expectedResponse   string
	}{
		{
			name:               "Invalid request body",
			inputBody:          "invalid-json",
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   `{"message": "invalid request body"}`,
		},
		{
			name: "Invalid email",
			inputBody: map[string]string{
				"email":    "invalid-email",
				"password": "StrongPass123!",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "invalid email address",
		},
		{
			name: "Invalid password",
			inputBody: map[string]string{
				"email":    "test@example.com",
				"password": "weak",
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "password does not meet complexity requirements",
		},
		{
			name: "User already exists",
			inputBody: map[string]string{
				"email":    "test@example.com",
				"password": "StrongPass123!",
			},
			mockResponse:       &pq.Error{Code: "23505"},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse:   "error registering user: user already exists",
		},
		{
			name: "Internal server error",
			inputBody: map[string]string{
				"email":    "test@example.com",
				"password": "StrongPass123!",
			},
			mockResponse:       errors.New("some error"),
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse:   "error registering user: database error",
		},
		{
			name: "Successful registration",
			inputBody: map[string]string{
				"email":    "test@example.com",
				"password": "StrongPass123!",
			},
			expectedStatusCode: http.StatusCreated,
			expectedResponse:   `{"message": "user registered successfully"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := &application.MockUserService{
				RegisterUserFunc: func(user domain.CreateUserRequest) error {
					return tt.mockResponse
				},
			}

			handler := NewUserHTTPHandler(mockService)

			body, _ := json.Marshal(tt.inputBody)
			req := httptest.NewRequest(http.MethodPost, "/users/register", bytes.NewReader(body))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			handler.RegisterUser(rec, req)

			if tt.expectedStatusCode != rec.Code {
				t.Errorf("expected status code %v, got %v", tt.expectedStatusCode, rec.Code)
			}
			if reflect.DeepEqual(rec.Body.Bytes(), []byte(tt.expectedResponse)) {
				t.Errorf("expected response body %s, got %s", tt.expectedResponse, rec.Body.String())
			}
		})
	}
}

// func TestUpdateUser(t *testing.T) {
// 	type args struct {
// 		expectedStatus int
// 		expectedBody   string
// 		userId         int
// 		requestBody    map[string]string
// 	}

// 	type fields struct {
// 		mockService *application.MockUserService
// 	}
// 	// Define test cases
// 	tests := []struct {
// 		name string
// 		fields
// 		args args
// 	}{
// 		{
// 			name: "Successful Update",
// 			args: args{
// 				requestBody: map[string]string{
// 					"firstName": "Updated",
// 					"lastName":  "Name",
// 				},
// 				expectedStatus: http.StatusOK,
// 				expectedBody:   `{"message":"User updated successfully"}`,
// 				userId:         1,
// 			},
// 			fields: fields{
// 				mockService: &application.MockUserService{
// 					UpdateUserDataFunc: func(user domain.User) error {
// 						return nil
// 					},
// 				},
// 			},
// 		},
// 		// {
// 		// 	name:        "Validation Error - Missing Fields",
// 		// 	requestBody: map[string]string{},
// 		// 	mockService: &application.MockUserService{
// 		// 		UpdateUserDataFunc: func(user domain.User) error {
// 		// 			return nil
// 		// 		},
// 		// 	},
// 		// 	expectedStatus: http.StatusBadRequest,
// 		// 	expectedBody:   `{"message":"First name and last name are required"}`,
// 		// 	userId:         1,
// 		// },
// 		// {
// 		// 	name: "Service Error",
// 		// 	requestBody: map[string]string{
// 		// 		"firstName": "Error",
// 		// 		"lastName":  "Case",
// 		// 	},
// 		// 	mockService: &application.MockUserService{
// 		// 		UpdateUserDataFunc: func(user domain.User) error {
// 		// 			return errors.New("failed to update user")
// 		// 		},
// 		// 	},
// 		// 	expectedStatus: http.StatusInternalServerError,
// 		// 	expectedBody:   `{"message":"Failed to update user"}`,
// 		// },
// 		{
// 			name: "missing user id",
// 			args: args{

// 				requestBody: map[string]string{
// 					"firstName": "Error",
// 					"lastName":  "Case",
// 				},
// 				expectedStatus: http.StatusBadRequest,
// 				expectedBody:   `{"message": "invalid user ID"}`,
// 			},
// 			fields: fields{
// 				mockService: &application.MockUserService{},
// 			},
// 		},
// 	}

// 	// Iterate over test cases
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			// Prepare request body
// 			body, _ := json.Marshal(tt.args.requestBody)
// 			fmt.Println("userid: ", tt.args.userId)
// 			// Create request and response recorder
// 			req := httptest.NewRequest("PUT", fmt.Sprintf("/users/%v", tt.args.userId), bytes.NewReader(body))
// 			req.Header.Set("Content-Type", "application/json")
// 			rec := httptest.NewRecorder()
// 			fmt.Printf("req: %+v", req)
// 			// Create handler with mock service
// 			handler := &UserHTTPHandler{UserService: tt.mockService}

// 			// Call the UpdateUser handler
// 			handler.UpdateUser(rec, req)
// 			res := rec.Result().Request
// 			fmt.Printf("res:%+v", res)
// 			// Assert status code
// 			if rec.Code != tt.args.expectedStatus {
// 				t.Errorf("expected status %d, got %d", tt.args.expectedStatus, rec.Code)
// 			}
// 			// Assert response body
// 			if reflect.DeepEqual(string(rec.Body.Bytes()), []byte(tt.args.expectedBody)) {
// 				t.Errorf("expected response body %s, got %s", tt.args.expectedBody, rec.Body.String())
// 			}
// 		})
// 	}
// }

// TestGetUserProfile_Success tests retrieving the user profile with valid authorization
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
