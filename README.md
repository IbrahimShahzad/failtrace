# failtrace
log trace on fail

## Use Case

- Accumulates logs and writes only if an error occured. If no error occurs, clears the buffer.
- Writes all logs whether Debug / Info / Warn or Error when Error occurs
- Adds a uuid for the request logger, easier to follow the request through the logs.
- After flush, new request-logger would have a different uuid.


### Example Output

```sh
> go run main.go
[a76c964f-83a2-4116-ad70-55cfc029d353] D: handling request
[a76c964f-83a2-4116-ad70-55cfc029d353] D: inside a
[a76c964f-83a2-4116-ad70-55cfc029d353] D: inside b
[a76c964f-83a2-4116-ad70-55cfc029d353] E: an error occurred in b
```
- for more see examples folder

> [!CAUTION]
> A single request-logger should not be used in between mutliple go-routines.


## Usage

```go
package main

import (
    "context"

    "github.com/IbrahimShahzad/failtrace/logger"
)

func main() {
    ctx := logger.WithLogger(context.Background())
    handle(ctx)
}

func handle(ctx context.Context) {
    log := logger.FromContext(ctx)
    log.Debug("handling request")
    err := someF(); err != nil {
        log.FlushIf(err)
    }

    log.Debug("everything is good")
    log.FlushIf(nil)
}

func someF() error {
    return errors.New("im an error")
}
```

The above should give you some thing like this:

```sh
[a76c964f-83a2-4116-ad70-55cfc029d353] D: handling request
[a76c964f-83a2-4116-ad70-55cfc029d353] E: im an error
```

> Copilot used for writing tests and benchmarks