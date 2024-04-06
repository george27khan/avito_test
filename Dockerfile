FROM golang:1.22-rc-alpine
WORKDIR /app
COPY go.mod ./

RUN go mod download
COPY . ./

RUN go build -o /avito_test

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 63342 # 8080
# command to be used to execute when the image is used to start a container
CMD [ "/avito_test" ]