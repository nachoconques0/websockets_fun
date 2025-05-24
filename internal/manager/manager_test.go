package manager_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"github.com/nachoconques0/websockets_fun/internal/manager"
	"github.com/nachoconques0/websockets_fun/internal/mocks"
)

func TestManager_ServeWS(t *testing.T) {
	tests := []struct {
		name          string
		message       string
		expectPublish bool
		mockReturnErr error
	}{
		{
			name:          "success - message published",
			message:       "hello world",
			expectPublish: true,
			mockReturnErr: nil,
		},
		{
			name:          "failure - publisher returns error",
			message:       "bad message",
			expectPublish: true,
			mockReturnErr: errors.New("redis publish failed"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockPublisher := mocks.NewMockPublisher(ctrl)

			if tc.expectPublish {
				mockPublisher.
					EXPECT().
					PublishMessage(tc.message).
					Return(tc.mockReturnErr).
					Times(1)
			}

			m := manager.NewManager(mockPublisher)

			server := httptest.NewServer(http.HandlerFunc(m.ServeWS))
			defer server.Close()

			wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

			conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
			assert.NoError(t, err)

			defer conn.Close()

			err = conn.WriteMessage(websocket.TextMessage, []byte(tc.message))
			assert.NoError(t, err)

		})
	}
}
