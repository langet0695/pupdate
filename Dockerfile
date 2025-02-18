# Use alpine build for pi-zero-w on linux/arm/v6 platform
FROM golang:1.24rc3-alpine3.20

WORKDIR /app

COPY go.mod go.sum ./
RUN ls
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /pupdate ./cmd/pupdate

EXPOSE 8080

RUN chmod +x /pupdate

CMD [ "/pupdate" ]