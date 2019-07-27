# matrix-key-server

Implementation of a key server for Matrix.

Support room: [#matrix-key-server:t2bot.io](https://matrix.to/#/#matrix-key-server:t2bot.io)

## Building and running

This project uses Go modules and requires Go 1.12 or higher. To enable modules, set `GO111MODULE=on`.

```bash
# Build
git clone https://github.com/t2bot/matrix-key-server.git
cd matrix-key-server
go build -v -o bin/matrix-key-server

# Run
./bin/matrix-key-server -address="0.0.0.0" -port=8080 -domain="keys.t2host.io" -postgres="postgres://username:password@localhost/dbname?sslmode=disable"
```

## Docker

```bash
docker run -it --rm -e "ADDRESS=0.0.0.0" -e "PORT=8080" -e "DOMAIN=keys.t2host.io" -e "POSTGRES=postgres://username:password@localhost/dbname?sslmode=disable" t2bot/matrix-key-server
```

Build your own by checking out the repository and running `docker build -t t2bot/matrix-key-server .`
