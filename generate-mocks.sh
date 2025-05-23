#!/bin/bash
mockgen --source=internal/manager/manager.go --destination=internal/mocks/mock_manager.go --package=mocks --mock_names=Publisher=MockPublisher
mockgen --source=internal/broadcaster/broadcaster.go --destination=internal/mocks/mock_broadcaster.go --package=mocks --mock_names=Broadcaster=MockBroadcaster
