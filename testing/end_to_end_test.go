package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"testing"
)

const (
	host = "http://localhost:8080/v1/users"
)

var (
	client = &http.Client{}
)

// TestEndToEndScenarioForTwoUsers tests a simple scenario for two persons' hypothetical viewing habits
func TestEndToEndScenarioForTwoUsers(t *testing.T) {
	as := assert.New(t)

	// becky starts watching
	as.Equal(http.StatusCreated, watchStream("becky", "rugby7"))
	as.Equal(http.StatusCreated, watchStream("becky", "tennis2"))

	as.Equal(http.StatusCreated, watchStream("charles", "boxing1"))

	// idempotency check
	as.Equal(http.StatusCreated, watchStream("charles", "cycling2"))
	as.Equal(http.StatusCreated, watchStream("charles", "cycling2"))
	as.Equal([]string{"boxing1", "cycling2"}, streamsWatched("charles"))

	// charles tries to watch karate and golf but exceeds his quota for the golf
	as.Equal(http.StatusCreated, watchStream("charles", "karate3"))
	as.Equal(http.StatusBadRequest, watchStream("charles", "golf4"))
	as.Equal([]string{"boxing1", "cycling2", "karate3"}, streamsWatched("charles"))

	// charles turns off boxing to watch golf
	as.Equal(http.StatusOK, finishStream("charles", "boxing1"))
	as.Equal(http.StatusCreated, watchStream("charles", "golf4"))
	as.Equal([]string{"cycling2", "golf4", "karate3"}, streamsWatched("charles"))

	// charles stops watching
	as.Equal(http.StatusOK, finishStream("charles", "cycling2"))
	as.Equal(http.StatusOK, finishStream("charles", "golf4"))
	as.Equal(http.StatusOK, finishStream("charles", "karate3"))
	as.Equal([]string{""}, streamsWatched("charles"))

	// becky stops watching
	as.Equal([]string{"rugby7", "tennis2"}, streamsWatched("becky"))
	as.Equal(http.StatusOK, finishStream("becky", "rugby7"))
	as.Equal(http.StatusOK, finishStream("becky", "tennis2"))
	as.Equal([]string{""}, streamsWatched("becky"))
}

func watchStream(userID, streamID string) int {
	req, err := http.NewRequest(
		"PUT",
		fmt.Sprintf("%v/%v/streams/%v", host, userID, streamID),
		nil,
	)
	if err != nil {
		panic(err)
	}
	status, _ := call(req)
	return status
}

func finishStream(userID, streamID string) int {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf("%v/%v/streams/%v", host, userID, streamID),
		nil,
	)
	if err != nil {
		panic(err)
	}
	status, _ := call(req)
	return status
}

func streamsWatched(userID string) []string {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("%v/%v", host, userID),
		nil,
	)
	if err != nil {
		panic(err)
	}
	_, body := call(req)
	streamIDs := strings.Split(body, ",")
	sort.Strings(streamIDs)
	return streamIDs
}

func call(req *http.Request) (int, string) {
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	data, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	return resp.StatusCode, string(data)
}
