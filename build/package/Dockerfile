FROM golang:alpine AS build

WORKDIR /app
RUN apk add --no-cache make

# Install dependencies.
COPY go.* ./
RUN go mod download

# Copy source files and build application.
COPY . /app
ARG VERSION
RUN VERSION=${VERSION} make -B && mv bin/* /app/bin/netadm

# Create minimal container image.
FROM alpine AS run
COPY --from=build /app/bin/netadm /netadm
CMD [ "/netadm" ]
