GOOS=linux GOARCH=amd64 go build -o bot ./main.go
docker build -t telegrambot .