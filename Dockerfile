FROM golang:1.22-rc-alpine as build
WORKDIR /app
COPY . ./

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/app ./cmd/app

FROM scratch

WORKDIR /app

COPY --from=build /app/bin/app ./bin/
COPY --from=build /app/.env ./.env

ENTRYPOINT ["./bin/app"]
# tells Docker that the container listens on specified network ports at runtime
#EXPOSE 8080
# command to be used to execute when the image is used to start a container
#CMD [ "/avito_test/cmd/app" ]