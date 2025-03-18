runtestserver:
	cd server && go run cmd/main.go

runstresstest:
	cd stress-tester && go run cmd/main.go

test:
	cd stress-tester && go test ./...

build:
	cd stress-tester && docker build -t stress-tester .
	cd server && docker build -t server .

run-server:
	docker run --name server -p 8080:8080 server