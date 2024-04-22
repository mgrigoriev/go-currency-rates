# Currency Converter

An application in Go for currency conversion using the Central Bank of Russia exchange rates.

### How It Works

* A separate goroutine downloads XML with currency rates every hour and caches them in memory.
* The HTTP server listens on port 8888 and provides two endpoints:
    * `/` — usage examples
    * `/convert` — currency converter

### Running the app (two ways)

1. Wrapper for `go run`:
    * `make run`
2. Building and running a Docker container:
    * `make d`

### Examples of API Requests

* 100 RUB → USD:
    - http://127.0.0.1:8888/convert/?amount=100&from=RUB&to=USD


* 100 EUR → RUB:
    - http://127.0.0.1:8888/convert/?amount=100&from=RUB&to=USD


* 1100 USD → AMD:
    - http://127.0.0.1:8888/convert/?amount=1100&from=USD&to=AMD

### Running the tests
`make test`