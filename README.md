# matrix-key-server

Implementation of a key server for Matrix.

Support room: [#matrix-key-server:t2bot.io](https://matrix.to/#/#matrix-key-server:t2bot.io)

**Caution**: Although this has notary server functionality, it is not yet recommended to point Synapse at this. It has not been tested - use at your own risk.

Demo: https://federationtester.matrix.org#keys.t2host.io

## Building and running

The key server will automatically generate itself a key to use on startup. The process is meant to be run 
only attached to a postgres instance and does not have any on-disk requirements other than the executable 
itself.

This project uses Go modules and requires Go 1.12 or higher. To enable modules, set `GO111MODULE=on`.

```bash
# Build
git clone https://github.com/t2bot/matrix-key-server.git
cd matrix-key-server
go build -v -o bin/matrix-key-server

# Run
./bin/matrix-key-server -address="0.0.0.0" -port=8080 -domain="keys.t2host.io" -postgres="postgres://username:password@localhost/dbname?sslmode=disable"
```

#### Docker

```bash
docker run -it --rm -e "ADDRESS=0.0.0.0" -e "PORT=8080" -e "DOMAIN=keys.t2host.io" -e "POSTGRES=postgres://username:password@localhost/dbname?sslmode=disable" t2bot/matrix-key-server
```

Build your own by checking out the repository and running `docker build -t t2bot/matrix-key-server .`

## Custom APIs

The key server exposes some custom APIs which may aide the development of homeservers or Matrix services.

#### `POST /_matrix/key/unstable/check_auth`

Verifies an auth header according to the Matrix specification. The `Authorization` header is passed through
and the remaining headers shown here demonstrate the additional information the key server needs. The content
for the API call is sent as the request body to this call.

**Caution**: Trusting this endpoint can be bad if you don't trust the key server. You should always do your own
auth wherever possible.

**Example request**:
```
POST /_matrix/key/unstable/check_auth
Authorization: X-Matrix origin="example.org",key="ed25519:auto",sig="ABCDEF..."
X-Keys-Method: GET
X-Keys-URI: /_matrix/federation/v1/publicRooms?include_all_networks=false&limit=20
X-Keys-Destination: dest.example.org

{... request body ...}
```

If the response is a `200 OK`, the server is authorized. All other responses should be considered unauthorized.
