# dbtest

HTTP testing made easy for layered Web applications in Go.


## Layered Web Applications

Non-trivial applications should be [layered][1].

![layered-app](layered-app.png)

For layered Web applications:

- HTTP Server/Client should be separated from Service (the business rules)
- HTTP Server should depend on an interface designed by Service
- HTTP Client should implement an interface designed by Service


## Installation

Make a custom build of [protogo](https://github.com/protogodev/protogo):

```bash
$ protogo build --plugin=github.com/protogodev/httptest
```

<details open>
  <summary> Usage </summary>

```bash
$ protogo httptest -h
Usage: protogo httptest --mode=STRING --spec=STRING <source-file> <interface-name>

Arguments:
  <source-file>       source-file
  <interface-name>    interface-name

Flags:
  -h, --help           Show context-sensitive help.

      --mode=STRING    generation mode (server or client)
      --spec=STRING    the test specification in YAML
      --out=STRING     output filename (default "./<srcPkgName>_<mode>_test.go")
      --fmt            whether to make the test code formatted
```
</details>


## Examples

See [examples/usersvc](examples/usersvc).


## Documentation

Check out the [Godoc][2].


## License

[MIT](LICENSE)


[1]: https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html
[2]: https://pkg.go.dev/github.com/protogodev/httptest
