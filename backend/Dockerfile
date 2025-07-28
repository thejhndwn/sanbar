FROM golang:1.24.5-bullseye

WORKDIR /app

COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /godocker ./cmd/main.go

EXPOSE 8080

CMD [ "/godocker" ]



