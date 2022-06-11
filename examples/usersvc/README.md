# User Store

This example illustrates how to apply declarative testing for a typical HTTP application.


## Run the server/client

Run the server:

```bash
$ go run cmd/server/main.go
2022/06/11 22:47:01 transport=HTTP addr=:8080
```

Run the client:

```bash
$ go run cmd/client/main.go
2022/06/11 22:47:03 GetUser ok: &{Name:foo Sex:male Birth:2022-01-01 00:00:00 +0000 UTC}
2022/06/11 22:47:03 ListUsers ok: [&{Name:foo Sex:female Birth:2022-01-01 00:00:00 +0000 UTC}]
```


## Prerequisites

1. Write the test specification for the HTTP server

    See [httptest.server.yaml](httptest.server.yaml).

2. Write the test specification for the HTTP client

   See [httptest.client.yaml](httptest.client.yaml).


## Generate tests

```bash
$ go generate
```

Generated tests:

- [HTTP server tests](usersvc_server_test.go)
- [HTTP client tests](usersvc_client.go)


## Run tests

```bash
$ go test -v -race
```

<details>
  <summary> Result </summary>

```bash
=== RUN   TestHTTPClient_GetUser
=== RUN   TestHTTPClient_GetUser/ok
=== RUN   TestHTTPClient_GetUser/err
--- PASS: TestHTTPClient_GetUser (0.00s)
    --- PASS: TestHTTPClient_GetUser/ok (0.00s)
    --- PASS: TestHTTPClient_GetUser/err (0.00s)
=== RUN   TestHTTPClient_ListUsers
=== RUN   TestHTTPClient_ListUsers/ok
--- PASS: TestHTTPClient_ListUsers (0.00s)
    --- PASS: TestHTTPClient_ListUsers/ok (0.00s)
=== RUN   TestHTTPClient_CreateUser
=== RUN   TestHTTPClient_CreateUser/ok
=== RUN   TestHTTPClient_CreateUser/err
--- PASS: TestHTTPClient_CreateUser (0.00s)
    --- PASS: TestHTTPClient_CreateUser/ok (0.00s)
    --- PASS: TestHTTPClient_CreateUser/err (0.00s)
=== RUN   TestHTTPClient_UpdateUser
=== RUN   TestHTTPClient_UpdateUser/ok
=== RUN   TestHTTPClient_UpdateUser/err
--- PASS: TestHTTPClient_UpdateUser (0.00s)
    --- PASS: TestHTTPClient_UpdateUser/ok (0.00s)
    --- PASS: TestHTTPClient_UpdateUser/err (0.00s)
=== RUN   TestHTTPClient_DeleteUser
=== RUN   TestHTTPClient_DeleteUser/ok
=== RUN   TestHTTPClient_DeleteUser/err
--- PASS: TestHTTPClient_DeleteUser (0.00s)
    --- PASS: TestHTTPClient_DeleteUser/ok (0.00s)
    --- PASS: TestHTTPClient_DeleteUser/err (0.00s)
=== RUN   TestHTTPServer_GetUser
=== RUN   TestHTTPServer_GetUser/ok
=== RUN   TestHTTPServer_GetUser/err
--- PASS: TestHTTPServer_GetUser (0.00s)
    --- PASS: TestHTTPServer_GetUser/ok (0.00s)
    --- PASS: TestHTTPServer_GetUser/err (0.00s)
=== RUN   TestHTTPServer_ListUsers
=== RUN   TestHTTPServer_ListUsers/ok
--- PASS: TestHTTPServer_ListUsers (0.00s)
    --- PASS: TestHTTPServer_ListUsers/ok (0.00s)
=== RUN   TestHTTPServer_CreateUser
=== RUN   TestHTTPServer_CreateUser/ok
=== RUN   TestHTTPServer_CreateUser/err
--- PASS: TestHTTPServer_CreateUser (0.00s)
    --- PASS: TestHTTPServer_CreateUser/ok (0.00s)
    --- PASS: TestHTTPServer_CreateUser/err (0.00s)
=== RUN   TestHTTPServer_UpdateUser
=== RUN   TestHTTPServer_UpdateUser/ok
=== RUN   TestHTTPServer_UpdateUser/err
--- PASS: TestHTTPServer_UpdateUser (0.00s)
    --- PASS: TestHTTPServer_UpdateUser/ok (0.00s)
    --- PASS: TestHTTPServer_UpdateUser/err (0.00s)
=== RUN   TestHTTPServer_DeleteUser
=== RUN   TestHTTPServer_DeleteUser/ok
=== RUN   TestHTTPServer_DeleteUser/err
--- PASS: TestHTTPServer_DeleteUser (0.00s)
    --- PASS: TestHTTPServer_DeleteUser/ok (0.00s)
    --- PASS: TestHTTPServer_DeleteUser/err (0.00s)
PASS
ok      github.com/protogodev/httptest/examples/usersvc 0.043s
```

</details>
