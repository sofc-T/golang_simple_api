build:
	@go build -o bin/task-manager ./Delivery/main.go

test:
	@go test $(go list ./... | grep -v '/repository/user_repo_test' | grep -v '/repository/task_repo_test') -v
