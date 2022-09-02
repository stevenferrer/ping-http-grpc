# ping-http-grpc

Ping via HTTP and [gRPC](https://grpc.io) on the same port using [cmux](https://github.com/soheilhy/cmux).

## Running the examples

Run the server.

```sh
$ go run cmd/server/main.go
2022/09/02 21:49:42 grpc server started.
2022/09/02 21:49:42 http server started.
2022/09/02 21:49:42 cmux started.
```

Ping via HTTP.

```sh
$ curl http://localhost:8080
pong
```

Ping gRPC.

```sh
$ go run cmd/pingclient/main.go
pong
```

## License

[MIT](./LICENSE)
