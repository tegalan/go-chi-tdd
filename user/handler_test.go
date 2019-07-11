package user

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/suite"
)

type UserHandlerTestSuite struct {
	suite.Suite
	Handler Handler
}

// SetupTest run before test
func (s *UserHandlerTestSuite) SetupTest() {

}

func (s *UserHandlerTestSuite) TestUserSignUp() {
	user := User{
		Name: "Moana",
	}
	st := new(MockStore)
	st.On("Create", &user).Return(nil)

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

	s.Equal(http.StatusOK, rr.Code, fmt.Sprintf("Unexpected code %v got body %s", rr.Code, rr.Body.String()))

	var body User
	merr := json.Unmarshal(rr.Body.Bytes(), &body)
	s.NoError(merr, fmt.Sprintf("Body got: %s", rr.Body.String()))

	s.Equal(user.Name, body.Name)
}

func (s *UserHandlerTestSuite) TestUserLogin() {
	pld, _ := json.Marshal(map[string]string{"email": "moana@motunui.is", "password": "heiheithechicken"})
	req, err := http.NewRequest("POST", "/user/login", bytes.NewBuffer(pld))
	req.Header.Set("Content-Type", "application/json")

	s.NoError(err)

	// Mock Handler store
	st := new(MockStore)
	moana := User{Name: "Moana", Email: "moana@motunui.is", ID: 1}
	st.On("FindByLogin", "moana@motunui.is", "heiheithechicken").Return(moana, nil)
	handler := Handler{
		store: st,
	}

	rr := httptest.NewRecorder()
	c := http.HandlerFunc(handler.Login)
	c.ServeHTTP(rr, req)

	s.Equal(http.StatusOK, rr.Code)

	var res User
	if err := json.Unmarshal(rr.Body.Bytes(), &res); err != nil {
		s.NoError(err)
	}

	s.Equal(moana.Name, res.Name, fmt.Sprintf("Unexpected result: %s expected %s", res.Name, moana.Name))
}

func TestUserHandlerSuite(t *testing.T) {
	suite.Run(t, new(UserHandlerTestSuite))
}
