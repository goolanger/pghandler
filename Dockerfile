# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM postgres:9-alpine

ENV PG_PORT 5432
ENV PG_USER postgres
ENV PG_OPERATION backup
ENV BACKUP_DIR /pgbackup
ENV DB_ENV prod

COPY backup.sh /usr/bin/backup.sh
COPY restore.sh /usr/bin/restore.sh
COPY fill.sh /usr/bin/fill.sh
COPY setpwd.sh /usr/bin/setpwd.sh

RUN chmod +x /usr/bin/backup.sh && chmod +x /usr/bin/restore.sh && chmod +x /usr/bin/fill.sh && chmod +x /usr/bin/setpwd.sh

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /app/main .

RUN mkdir -p ${BACKUP_DIR}

# Environment Variables
ENV LOG_FILE_LOCATION=${BACKUP_DIR}/pghandler.log

# Expose port 8080 to the outside world
EXPOSE 8091

# Declare volumes to mount
VOLUME [${BACKUP_DIR}]

# Command to run the executable
CMD ["./main"]