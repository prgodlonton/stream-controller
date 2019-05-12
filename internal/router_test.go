package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/pgodlonton/stream-controller/testing/mocks"
	"net/http"
	"testing"
)

func TestShouldReturnCreatedWhenStreamIsAdded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().Add("alan", "boxing1").MinTimes(1).Return(nil)
	store.EXPECT().Remove(gomock.Any(), gomock.Any()).MaxTimes(0)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusCreated)

	r := createHTTPRequest("PUT", "v1/users/alan/streams/boxing1")

	router := NewRouter(store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnBadRequestWhenUsersHasReachedStreamQuotaLimit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().Add("michelangelo", "bobsleigh32").MinTimes(1).Return(exceededStreamsQuota)
	store.EXPECT().Remove(gomock.Any(), gomock.Any()).MaxTimes(0)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusBadRequest)

	r := createHTTPRequest("PUT", "v1/users/michelangelo/streams/bobsleigh32")

	router := NewRouter(store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnInternalServerErrorWhenStreamCannotBeAdded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().Add("bob", "tennis2").MinTimes(1).Return(errors.New("intentional error"))
	store.EXPECT().Remove(gomock.Any(), gomock.Any()).MaxTimes(0)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusInternalServerError)

	r := createHTTPRequest("PUT", "v1/users/bob/streams/tennis2")

	router := NewRouter(store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnActiveStreamsForUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().Get("cassandra").MaxTimes(1).Return(
		[]string{"boxing16", "tennis42", "sumo89"},
		nil,
	)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().Write([]byte("boxing16,tennis42,sumo89"))

	r := createHTTPRequest("GET", "v1/users/cassandra")

	router := NewRouter(store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnInternalServerErrorWhenActiveStreamsCannotBeRead(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().Get("rachel").MaxTimes(1).Return(
		[]string{},
		errors.New("intentional error"),
	)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusInternalServerError)

	r := createHTTPRequest("GET", "v1/users/rachel")

	router := NewRouter(store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnOKWhenStreamIsRemoved(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().Add(gomock.Any(), gomock.Any()).MaxTimes(0)
	store.EXPECT().Remove("charlie", "snooker3").MinTimes(1).Return(nil)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusOK)

	r := createHTTPRequest("DELETE", "v1/users/charlie/streams/snooker3")

	router := NewRouter(store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnInternalServerErrorWhenStreamCannotBeRemoved(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().Add(gomock.Any(), gomock.Any()).MaxTimes(0)
	store.EXPECT().Remove("duncan", "nfl4").MinTimes(1).Return(errors.New("intentional error"))

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusInternalServerError)

	r := createHTTPRequest("DELETE", "v1/users/duncan/streams/nfl4")

	router := NewRouter(store)
	router.ServeHTTP(w, r)
}

func createHTTPRequest(method, url string) *http.Request {
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("http://localhost:8080/%v", url),
		bytes.NewReader([]byte{}),
	)
	if err != nil {
		panic(err)
	}
	return req
}
