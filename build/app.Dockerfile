FROM golang:1.22.5-alpine3.20 as build
WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download && go mod verify
COPY cmd /src/cmd
COPY internal /src/internal
RUN go build -o /app/app /src/cmd/golang_task/golang_task.go

FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/app /app/app
EXPOSE 88
CMD ["/app/app"]