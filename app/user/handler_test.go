package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	Router *chi.Mux
}

// SetupTest run before test
func (s *UserHandlerTestSuite) SetupTest() {
	s.Router = chi.NewMux()
}

func (s *UserHandlerTestSuite) TestGetRoutes() {
	h := Handler{}
	h.RegisterRouter(s.Router)
}

func (s *UserHandlerTestSuite) TestUserSignUp() {
	tests := []struct {
		user    User
		payload render.M
		store   error
		code    int
	}{
		{
			payload: render.M{},
			user:    User{},
			store:   nil,
			code:    http.StatusUnprocessableEntity,
		},
		{
			payload: render.M{"name": "", "": "", "password": ""},
			user:    User{},
			store:   nil,
			code:    http.StatusUnprocessableEntity,
		},
		{
			payload: render.M{"name": "Moana", "email": "moana@motunui.is", "password": "heiheithechicken"},
			user:    User{Name: "Moana", Email: "moana@motunui.is", Password: "heiheithechicken"},
			store:   nil,
			code:    http.StatusOK,
		},
		{
			payload: render.M{"name": "Tekka", "email": "tekka@motunui.is", "password": "tefiti"},
			user:    User{Name: "Tekka", Email: "tekka@motunui.is", Password: "tefiti"},
			store:   errors.New("Unable to save user"),
			code:    http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		user := test.user
		st := new(MockStore)
		st.On("Create", &user).Return(test.store)

		handler := NewHandler(st)

		pld, _ := json.Marshal(test.payload)
		req, err := http.NewRequest("POST", "/user", bytes.NewBuffer(pld))
		req.Header.Set("Content-Type", "application/json")

		s.NoError(err)

		rr := httptest.NewRecorder()
		client := http.HandlerFunc(handler.SignUp)
		client.ServeHTTP(rr, req)

		s.Equal(test.code, rr.Code, fmt.Sprintf("Unexpected code %v got body %s", rr.Code, rr.Body.String()))

		if rr.Code == http.StatusOK {
			var body User
			merr := json.Unmarshal(rr.Body.Bytes(), &body)
			s.NoError(merr, fmt.Sprintf("Body got: %s", rr.Body.String()))

			s.Equal(user.Name, body.Name)
		}
	}
}

func (s *UserHandlerTestSuite) TestUserLogin() {

	tests := []struct {
		email    string
		password string
		user     User
		expected int
	}{
		{ //Existing user
			email:    "moana@motunui.is",
			password: "heiheithechicken",
			user:     User{Name: "Moana", Email: "moana@motunui.is", Password: "$2y$12$LX9yhsMiB8n9wmAhRfQLRugtT3nVBhC.DIbdr3hJcYBI4CV5nUfIO"},
			expected: http.StatusOK,
		},
		{ // Empty payload
			user:     User{},
			expected: http.StatusUnprocessableEntity,
		},
		{ // User Not Exist
			email:    "tekka@motunui.is",
			password: "tefiti",
			user:     User{Name: "Tekka", Email: "tekka@motunui.is", Password: "tefiti"},
			expected: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {

		pld, _ := json.Marshal(map[string]string{"email": test.email, "password": test.password})
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(pld))
		req.Header.Set("Content-Type", "application/json")

		s.NoError(err)

		// Mock Handler store
		st := new(MockStore)
		st.On("FindByLogin", "moana@motunui.is", "heiheithechicken").Return(test.user, nil)
		st.On("FindByLogin", "tekka@motunui.is", "tefiti").Return(User{}, errors.New("User Not Found"))

		handler := NewHandler(st)

		rr := httptest.NewRecorder()
		c := http.HandlerFunc(handler.Login)
		c.ServeHTTP(rr, req)

		s.Equal(test.expected, rr.Code, rr.Body.String())

		var res LoginResponse
		if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
			s.NoError(err)
		}

		if rr.Code == http.StatusOK {
			s.Equal(test.user.Name, res.Name, fmt.Sprintf("Unexpected result: %s expected %s", res.Name, test.user.Name))
			s.NotEmpty(res.Token, "Token empty")
		}

	}
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
