#!/bin/bash
mockgen --source=internal/manager/controller/manager/controller.go --destination=internal/mocks/mock_manager_controller.go --package=mocks --mock_names=Service=MockService
mockgen --source=internal/manager/service/manager/service.go --destination=internal/mocks/mock_manager_service.go --package=mocks --mock_names=Publisher=MockPublisher

mockgen --source=internal/broadcaster/controller/subscriber/controller.go --destination=internal/mocks/mock_broadcaster_controller.go --package=mocks --mock_names=Service=MockBroadcasterService
