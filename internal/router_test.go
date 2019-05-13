package internal

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/prgodlonton/stream-controller/testing/mocks"
	"go.uber.org/zap"
	"net/http"
	"testing"
)

var (
	noopLogger = zap.NewNop().Sugar()
)

func TestShouldReturnCreatedWhenStreamIsAdded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().AddStream("alan", "boxing1").MinTimes(1).Return(nil)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusCreated)

	r := createHTTPRequest("PUT", "v1/users/alan/streams/boxing1")

	router := NewRouter(noopLogger, store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnBadRequestWhenUserHasReachedStreamQuotaLimit(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().AddStream("michelangelo", "bobsleigh32").MinTimes(1).Return(exceededStreamsQuota)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusBadRequest)

	r := createHTTPRequest("PUT", "v1/users/michelangelo/streams/bobsleigh32")

	router := NewRouter(noopLogger, store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnInternalServerErrorWhenStreamCannotBeAdded(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().AddStream("bob", "tennis2").MinTimes(1).Return(errors.New("intentional error"))

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusInternalServerError)

	r := createHTTPRequest("PUT", "v1/users/bob/streams/tennis2")

	router := NewRouter(noopLogger, store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnActiveStreamsForUser(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().GetStreams("cassandra").MaxTimes(1).Return(
		[]string{"boxing16", "tennis42", "sumo89"},
		nil,
	)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().Write([]byte("boxing16,tennis42,sumo89"))

	r := createHTTPRequest("GET", "v1/users/cassandra")

	router := NewRouter(noopLogger, store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnInternalServerErrorWhenActiveStreamsCannotBeRead(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().GetStreams("rachel").MaxTimes(1).Return(
		[]string{},
		errors.New("intentional error"),
	)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusInternalServerError)

	r := createHTTPRequest("GET", "v1/users/rachel")

	router := NewRouter(noopLogger, store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnInternalServerErrorWhenWritingStreamListToResponseWriterFails(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().GetStreams("rodney").MaxTimes(1).Return(
		[]string{"boxing16", "tennis42", "sumo89"},
		nil,
	)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().Write(gomock.Any()).Return(0, errors.New("intentional error"))
	w.EXPECT().WriteHeader(http.StatusInternalServerError)

	r := createHTTPRequest("GET", "v1/users/rodney")

	router := NewRouter(noopLogger, store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnOKWhenStreamIsRemoved(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().RemoveStream("charlie", "snooker3").MinTimes(1).Return(nil)

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusOK)

	r := createHTTPRequest("DELETE", "v1/users/charlie/streams/snooker3")

	router := NewRouter(noopLogger, store)
	router.ServeHTTP(w, r)
}

func TestShouldReturnInternalServerErrorWhenStreamCannotBeRemoved(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	store := mocks.NewMockStore(mockCtrl)
	store.EXPECT().RemoveStream("duncan", "nfl4").MinTimes(1).Return(errors.New("intentional error"))

	w := mocks.NewMockResponseWriter(mockCtrl)
	w.EXPECT().WriteHeader(http.StatusInternalServerError)

	r := createHTTPRequest("DELETE", "v1/users/duncan/streams/nfl4")

	router := NewRouter(noopLogger, store)
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
