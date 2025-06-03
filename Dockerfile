# Stage 1: Build Vue.js frontend
FROM node:18-alpine AS vue_builder

WORKDIR /app

# Copy package manager files
COPY frontend/package.json frontend/package-lock.json ./frontend/
# COPY frontend/pnpm-lock.yaml ./frontend/ # Uncomment if using pnpm
# COPY frontend/yarn.lock ./frontend/ # Uncomment if using yarn

# Install frontend dependencies
WORKDIR /app/frontend
RUN npm install
# RUN pnpm install # Uncomment if using pnpm
# RUN yarn install # Uncomment if using yarn

# Copy the rest of the frontend code
WORKDIR /app 
COPY frontend/ ./frontend/

# Build the frontend application
WORKDIR /app/frontend
RUN npm run build 
# RUN pnpm build # Uncomment if using pnpm
# RUN yarn build # Uncomment if using yarn

# Stage 2: Build Go backend and serve frontend
FROM golang:1.21-alpine AS go_runner
# Using 1.21-alpine as specified in original prompt, though current Dockerfile used 1.23

WORKDIR /app

# Set Go environment variables
ENV CGO_ENABLED=0 GOOS=linux

# Copy Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy the Go application source code
COPY . . 
# This includes main.go, internal/, etc. 
# It also copies the 'frontend' source directory again, which is not ideal but usually harmless.
# A more precise copy would be `COPY internal/ ./internal/` and `COPY main.go .` etc.

# Build the Go application
# Output will be in /app/main
RUN go build -ldflags="-w -s" -o /app/main .

# Copy built Vue.js app from the vue_builder stage
# The Go app will serve files from /app/frontend/dist
COPY --from=vue_builder /app/frontend/dist /app/frontend/dist

# Expose the application port
EXPOSE 8080

# Set the command to run the Go application
# The Go application will serve the static files and the API
CMD ["/app/main"]