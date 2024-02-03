FROM golang:1.22-rc-alpine
WORKDIR /app
COPY go.mod ./

RUN go mod download
COPY . ./

RUN go build -o /discord_bot

# tells Docker that the container listens on specified network ports at runtime
EXPOSE 8080
# command to be used to execute when the image is used to start a container
CMD [ "/discord_bot" ]