package service_test

import (
	"testing"

	"github.com/nachoconques0/websockets_fun/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	service "github.com/nachoconques0/websockets_fun/internal/manager/service/manager"
)

func TestService_PublishMessage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockPublisher := mocks.NewMockPublisher(ctrl)
	mockPublisher.EXPECT().
		PublishMessage("test").
		Return(nil).
		Times(1)

	s := service.New(mockPublisher)

	err := s.PublishMessage("test")
	assert.NoError(t, err)
}
