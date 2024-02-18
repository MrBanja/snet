# snet

`snet` is a simple Go library that provides helper functions for the `net/http` package. It simplifies the process of creating HTTP requests, (un)marshalling, constructing URLs, and gracefully shutting down servers.

## Installation

To install `snet`, use the following command:

```bash
go get github.com/mrbanja/snet/v2
```

## Usage

Here's a brief overview of the functions provided by `snet`:

### NewRequest

`NewRequest` creates a new HTTP request with the given method, URL, and body. The body is an optional JSON-serializable struct.

```go
type User struct {
    Name string `json:"name"`
    Age int `json:"age"`
}
var body = &User{Name: "John", Age: 30}
req, err := snet.NewRequest(context.TODO(), "POST", "https://example.com", body)
```

### UnmarshalResp (UnmarshalRequest, Unmarshal)

`Unmarshal` unmarshals the response body into a new instance of the provided type.

```go
type User struct {
    Name string `json:"name"`
    Age int `json:"age"`
}
user, err := sreq.UnmarshalResp[User](resp)
```

### U

`U` creates a new url.URL instance with the given path appended to the base URL.

```go
u, err := snet.U("https://example.com/api", "/user/create")
```

### ListenAndServe

`ListenAndServe` starts the server and listens for signals to shut down the server.

```go
err := snet.ListenAndServe(ctx, server, logger, os.Interrupt)
```

### Errors

`snet` also provides custom errors. You can:
- Create a new error using `NewWrongStatusError(...)`
- Check if an error is of a specific type using `IsWrongStatusError(...)`

#### WrongStatusError
This error type includes the response body, response status code, request URL, and request method.

```go
if req.StatusCode != http.StatusOK {
    return snet.NewWrongStatusError(req)
}
```
The output will be:
```
wrong status code for [GET https://example.com]; Code: [405] with reponse [Method Not Allowed, Use POST instead]
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

`snet` is released under the MIT License.
