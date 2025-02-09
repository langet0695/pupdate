FROM golang:1.23.5

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
# COPY cmd internal tmp ./
# COPY internal/* ./
# COPY src/.env ./
# COPY tmp ./tmp



COPY . .
# RUN go build -ldflags '-w -s' -a -o ./bin/api ./cmd/pupdate 

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /pupdate ./cmd/pupdate

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
# But we can (optionally) document in the Dockerfile what ports
# the application is going to listen on by default.
# https://docs.docker.com/engine/reference/builder/#expose
EXPOSE 8080
RUN chmod +x /pupdate
# Run
CMD [ "/pupdate" ]