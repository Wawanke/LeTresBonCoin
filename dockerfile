FROM golang:1.19 AS builder
WORKDIR /usr/src/

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
EXPOSE 3000:3000
RUN go build -o /godocker
CMD [ "/godocker" ]