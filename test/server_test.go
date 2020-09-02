package test

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	usercontr "addreality_t/internal/addreality/users"
	userdp "addreality_t/internal/addreality/users/withdb"
	"addreality_t/internal/resources"
	"addreality_t/internal/restapi"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// Create server with mock Redis connection
func createTestServer() restapi.ServiceHandlers {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	slogger := logger.Sugar()

	rsc, err := resources.New(slogger)
	if err != nil {
		slogger.Fatalw("Can't initialize resources.", "err", err)
	}
	defer func() {
		rsc.Release()
	}()

	userdb := userdp.New(rsc.DB)
	uc := usercontr.NewController(userdb)

	service := restapi.ServiceHandlers{
		UserController: uc,
	}
	return service
}

func makeRequest(t *testing.T, server restapi.ServiceHandlers, url string) *httptest.ResponseRecorder {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(server.CountUser()).ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}
	return rr
}

// Testing user handler
func TestUserCountHandler(t *testing.T) {
	server := createTestServer()
	rr := makeRequest(t, server, "/?user_id=abc")
	assert.Equal(t, rr.Code, 200, "By default user count. Handler count user working")
}

// Testing user handler
func TestCountHandler(t *testing.T) {
	server := createTestServer()

	rr := makeRequest(t, server, "/count")
	resp, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
	}
	byteToInt, _ := strconv.Atoi(string(resp))
	assert.Equal(t, byteToInt, 0, "By default no robots in service. Handler counter working")
}

// Testing One user in Robot
func TestUserCountRobot(t *testing.T) {
	server := createTestServer()
	for i := 0; i < 200; i++ {
		makeRequest(t, server, "/?user_id=abc")
	}

	req, err := http.NewRequest("GET", "/count", nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(server.GetRobotUser()).ServeHTTP(rr, req)

	//Confirm the response has the right status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}
	resp, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
	}
	byteToInt, _ := strconv.Atoi(string(resp))
	assert.Equal(t, byteToInt, 1, "One user is Robot")
}

// Testing several users in robot
func TestUsersCountRobot(t *testing.T) {
	server := createTestServer()
	for i := 0; i < 101; i++ {
		makeRequest(t, server, "/?user_id=abc")
	}

	for i := 0; i < 99; i++ {
		makeRequest(t, server, "/?user_id=abc11")
	}

	req, err := http.NewRequest("GET", "/count", nil)
	if err != nil {
		t.Errorf("An error occurred. %v", err)
	}

	rr := httptest.NewRecorder()
	http.HandlerFunc(server.GetRobotUser()).ServeHTTP(rr, req)

	//Confirm the response has the right status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Status code differs. Expected %d .\n Got %d instead", http.StatusOK, status)
	}
	resp, err := ioutil.ReadAll(rr.Body)
	if err != nil {
		log.Fatal(err)
	}
	byteToInt, _ := strconv.Atoi(string(resp))
	assert.Equal(t, byteToInt, 1, "One user is Robot. Second user is normal")
}
