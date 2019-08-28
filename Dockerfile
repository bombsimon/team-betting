FROM golang:alpine

# Set environment
ENV \
    DB_DSN="betting:betting@tcp(mysql:3306)/betting?charset=utf8&parseTime=True" \
    ADD_DATA=1


# Setup workdir
WORKDIR /app

# Copy repository
COPY ./ /app

# Install git to get go modules working
RUN apk add --update \
    bash \
    git

# Install wait-for-it.sh
RUN git clone https://github.com/vishnubob/wait-for-it.git

# Build migration tools
RUN go build -o gorm-migrate cmd/gorm-migrate/main.go

# Build betting application
RUN go build -o team-betting cmd/betting/main.go

# Start the application
CMD [ \
    "./wait-for-it/wait-for-it.sh", \
    "-t", "0", \
    "mysql:3306", \
    "--", \
    "sh", "-c", "./gorm-migrate && ./team-betting" \
]
