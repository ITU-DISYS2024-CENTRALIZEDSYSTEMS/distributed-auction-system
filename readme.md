# Distributed auction house demo
A distributed auction house using active replication

## Usage

### Make evnironment file of ports

Create a file named ``.env`` in the root directory with the wanted amount of ports.
```sh
PORTS=50051,50052,50053
```
_Each port has to be seperated with `,` after each port_

### Run the servers
Run the number of servers denoted by the amount of servers in the .env file.
```sh
go run server/server.go
```

### Run the clients

Run the amount of desired clients
```sh
go run client/client.go
```

Chose a username and start bidding or fetch the results with

```sh
bid <amount>
```

or

```sh
results
```

_Note:_ the auction starts as soon as the servers become online.
