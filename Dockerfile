FROM golang:1.22 as builder

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY main.go .
RUN go mod download
RUN go build -o ./app


FROM gcr.io/distroless/base-debian12
WORKDIR /appdir
COPY --from=builder /build/app ./app

CMD [ "/appdir/app" ]
