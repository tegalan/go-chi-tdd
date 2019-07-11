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
	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	Handler Handler
}

// SetupTest run before test
func (s *UserHandlerTestSuite) SetupTest() {

}

func (s *UserHandlerTestSuite) TestGetRoutes() {
	h := s.Handler.Routes()

	r := chi.NewRouter()
	r.Mount("/user", h)
}

func (s *UserHandlerTestSuite) TestUserSignUp() {
	tests := []struct {
		user  User
		store error
		code  int
	}{
		{
			user:  User{},
			store: nil,
			code:  http.StatusUnprocessableEntity,
		},
		{
			user:  User{Name: "Moana", Email: "moana@motunui.is", Password: "heiheithechicken"},
			store: nil,
			code:  http.StatusOK,
		},
		{
			user:  User{Name: "Tekka", Email: "tekka@motunui.is", Password: "tefiti"},
			store: errors.New("Unable to save user"),
			code:  http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		user := test.user
		st := new(MockStore)
		st.On("Create", &user).Return(test.store)

		handler := Handler{
			store: st,
		}

		pld, _ := json.Marshal(user)
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
		user     User
		expected int
	}{
		{
			user:     User{Name: "Moana", Email: "moana@motunui.is", Password: "heiheithechicken"},
			expected: http.StatusOK,
		},
		{
			user:     User{},
			expected: http.StatusUnprocessableEntity,
		},
		{
			user:     User{Name: "Tekka", Email: "tekka@motunui.is", Password: "tefiti"},
			expected: http.StatusUnauthorized,
		},
	}

	for _, test := range tests {

		pld, _ := json.Marshal(map[string]string{"email": test.user.Email, "password": test.user.Password})
		req, err := http.NewRequest("POST", "/", bytes.NewBuffer(pld))
		req.Header.Set("Content-Type", "application/json")

		s.NoError(err)

		// Mock Handler store
		st := new(MockStore)
		st.On("FindByLogin", "moana@motunui.is", "heiheithechicken").Return(test.user, nil)
		st.On("FindByLogin", "tekka@motunui.is", "tefiti").Return(User{}, errors.New("User Not Found"))
		handler := Handler{
			store: st,
		}

		rr := httptest.NewRecorder()
		c := http.HandlerFunc(handler.Login)
		c.ServeHTTP(rr, req)

		s.Equal(test.expected, rr.Code)

		var res User
		if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
			s.NoError(err)
		}

		if rr.Code == http.StatusOK {
			s.Equal(test.user.Name, res.Name, fmt.Sprintf("Unexpected result: %s expected %s", res.Name, test.user.Name))
		}

	}
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
