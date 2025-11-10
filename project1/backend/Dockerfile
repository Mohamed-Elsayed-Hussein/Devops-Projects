# ------ Builder stage ------ 
FROM golang:tip-alpine3.21  AS builder 

WORKDIR /app

# install dependencies and maintain caching 
COPY go.mod go.sum ./
RUN go mod download

# Copy the main code 
COPY main.go ./
RUN go build -ldflags="-s -w" -o backend



# ------ final stage ------ 
FROM  golang:tip-alpine3.21 
# @sha256:91b55682ac800beeac41d781ed70051ac5cdf1c5e47708af6ff9997e75936550

WORKDIR /app

COPY --from=builder /app/backend /app/backend

RUN adduser -DH GoApp
RUN chown -R GoApp:GoApp /app

EXPOSE 8000


USER GoApp
CMD ["./backend"]



