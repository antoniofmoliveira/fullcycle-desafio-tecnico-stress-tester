runtestserver:
	cd server && go run cmd/main.go

runstresstest:
	cd stress-tester && go run cmd/main.go

test:
	cd stress-tester && go test ./...