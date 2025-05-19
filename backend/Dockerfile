FROM golang:1.24.0
WORKDIR app/
COPY . .
RUN go mod tidy && go build -o main ./go_code/app.go