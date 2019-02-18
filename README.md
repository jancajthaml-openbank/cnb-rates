# cnb-rates

## Description

microservice to provide offline CNB fx rates access synchronizing daily rates
automatically.

## Compilation

### ARM

To compile for `ARM` use the following environment.

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=7 go build -v
```

### AMD64

To compile for `AMD64` use the following environment.

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v
```

## License

Licensed under Apache 2.0 see LICENSE.md for details

## Author

Jan Cajthaml (a.k.a johnny)
