package interfaces

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

func TestUpdateUser(t *testing.T) {
	tests := []struct {
		name            string
		PathValue       string
		body            string
		mockUserService application.UserServiceInterface
		expectedStatus  int
		expectedBody    string
	}{

		{
			name:            "Invalid user ID",
			PathValue:       "test",
			body:            "{}",
			mockUserService: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "{\"message\": \"invalid user ID\"}\n",
		},
		{
			name:            "Invalid request body",
			PathValue:       "1",
			body:            "invalid-json",
			mockUserService: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "{\"message\": \"invalid request body\"}\n",
		},
		{
			name:            "Invalid email",
			PathValue:       "1",
			body:            `{"email": "invalid-email"}`,
			mockUserService: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "invalid email format\n",
		},
		{
			name:            "Invalid password",
			PathValue:       "1",
			body:            `{"password": "short"}`,
			mockUserService: nil,
			expectedStatus:  http.StatusBadRequest,
			expectedBody:    "password must be at least 8 characters\n",
		},
		{
			name:      "Successful update",
			PathValue: "1",
			body:      `{"email": "valid@example.com", "password": "ValidPassword123"}`,
			mockUserService: &application.MockUserService{
				UpdateUserDataFunc: func(user domain.User) error {
					return nil
				},
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "{\"message\":\"user updated successfully\"}\n",
		},
		{
			name:      "Service error",
			PathValue: "1",
			body:      `{"email": "valid@example.com"}`,
			mockUserService: &application.MockUserService{
				UpdateUserDataFunc: func(user domain.User) error {
					return errors.New("database error")
				},
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "error updating user: database error\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			h := UserHTTPHandler{
				UserService: tt.mockUserService,
			}

			req := httptest.NewRequest(http.MethodPut, "/user/{id}", strings.NewReader(tt.body))
			req.SetPathValue("id", tt.PathValue)

			fmt.Printf("%+v", req)

			w := httptest.NewRecorder()

			h.UpdateUser(w, req)

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(res.Body)
			if buf.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, buf.String())
			}
		})
	}
}

func TestGetUserProfile(t *testing.T) {
	utils.JWTSecretKey = []byte("test secret")
	f := "Existing"
	l := "User"
	u := "John"
	mockUser := &domain.User{
		ID:        1,
		Username:  &u,
		FirstName: &f,
		LastName:  &l,
	}

	tests := []struct {
		name                string
		userIDInContext     int
		userIDFromURL       string
		expectedStatusCode  int
		mockGetUserByIDFunc func(id int) (*domain.User, error)
		expectedBody        interface{}
	}{
		{
			name:               "Unauthorized - missing user ID in context",
			userIDInContext:    0,
			userIDFromURL:      "1",
			expectedStatusCode: http.StatusForbidden,
			mockGetUserByIDFunc: func(id int) (*domain.User, error) {
				return nil, nil
			},
			expectedBody: "Unauthorized\n",
		},
		{
			name:               "Unauthorized - user ID mismatch",
			userIDInContext:    2,
			userIDFromURL:      "1",
			expectedStatusCode: http.StatusForbidden,
			mockGetUserByIDFunc: func(id int) (*domain.User, error) {
				return nil, nil
			},
			expectedBody: "Unauthorized\n",
		},
		{
			name:               "User not found",
			userIDInContext:    1, // Matches userIDFromURL
			userIDFromURL:      "1",
			expectedStatusCode: http.StatusNotFound,
			mockGetUserByIDFunc: func(id int) (*domain.User, error) {
				return nil, sql.ErrNoRows
			},
			expectedBody: "User not found\n",
		},

		{
			name:               "Internal server error",
			userIDInContext:    1,
			userIDFromURL:      "1",
			expectedStatusCode: http.StatusInternalServerError,
			mockGetUserByIDFunc: func(id int) (*domain.User, error) {
				return nil, errors.New("some internal error")
			},
			expectedBody: "Internal server error\n",
		},
		{
			name:               "Successful user profile retrieval",
			userIDInContext:    1,
			userIDFromURL:      "1",
			expectedStatusCode: http.StatusOK,
			mockGetUserByIDFunc: func(id int) (*domain.User, error) {
				return mockUser, nil
			},
			expectedBody: func() string {
				data, _ := json.Marshal(mockUser)
				return string(data) + "\n" // Add newline to match actual response
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up the mock service
			mockService := &application.MockUserService{
				GetUserByIDFunc: tt.mockGetUserByIDFunc,
			}

			// Create handler with mock service
			handler := NewUserHTTPHandler(mockService)

			// Create the request
			req := httptest.NewRequest(http.MethodGet, "/users/{id}/profile", nil)
			req = req.WithContext(context.WithValue(context.Background(), userIDKey, tt.userIDInContext))
			req = req.WithContext(context.WithValue(req.Context(), isAdminKey, false))
			req.SetPathValue("id", tt.userIDFromURL)

			// Create a ResponseRecorder to capture the response
			r := httptest.NewRecorder()

			// Call the handler
			handler.GetUserProfile(r, req)

			res := r.Result()
			defer res.Body.Close()

			// Check the status code
			if res.StatusCode != tt.expectedStatusCode {
				t.Errorf("expected status code %d, got %d", tt.expectedStatusCode, res.StatusCode)
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(res.Body)
			if buf.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, buf.String())
			}
		})
	}
}
