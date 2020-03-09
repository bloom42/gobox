<p align="center">
  <h3 align="center">graphql-go</h3>
  <p align="center">Easy to use GraphQL client for Go. Dependency free.</p>
</p>

[![GoDoc](https://godoc.org/gitlab.com/bloom42/libs/graphql-go?status.svg)](https://pkg.go.dev/gitlab.com/bloom42/libs/graphql-go)
[![CI/CD](https://gitlab.com/bloom42/libs/graphql-go/badges/master/pipeline.svg)](https://gitlab.com/bloom42/libs/graphql-go/pipelines)

1. [Installation](#configuration)
2. [Usage](#usage)
3. [Contributing](#contributing)
4. [Forum](#contributing)
5. [Security](#security)
6. [License](#license)

-------------------

Easy to use GraphQL client for Go.

* Simple, familiar API
* Respects `context.Context` timeouts and cancellation
* Build and execute any kind of GraphQL request
* Use strong Go types for response data
* Use variables and upload files
* Simple error handling

## Installation

Make sure you have a working Go environment. To install graphql, simply Run:

```
$ go get -u gitlab.com/bloom42/libs/graphql-go
```

## Usage

```go
import "context"

// create a client (safe to share across requests)
client := graphql.NewClient("https://bloom42.com/api/graphql")

// make a request
req := graphql.NewRequest(`
    query ($key: String!) {
        items (id:$key) {
            field1
            field2
            field3
        }
    }
`)

// set any variables
req.Var("key", "value")

// set header fields
req.Header.Set("Cache-Control", "no-cache")

// define a Context for the request
ctx := context.Background()

// Do it and capture the response
var respData ResponseStruct
if err := client.Do(ctx, req, &respData); err != nil {
    log.Fatal(err)
}
```

### File support via multipart form data

By default, the package will send a JSON body. To enable the sending of files, you can opt to
use multipart form data instead using the `UseMultipartForm` option when you create your `Client`:

```
client := graphql.NewClient("https://bloom42.com/api/graphql", graphql.UseMultipartForm())
```

## Contributing

Thank you for your interest in contributing! Please refer to
https://gitlab.com/bloom42/wiki/-/wikis/organization/contributing for guidance.


## Forum

When you want to talk about something that does not necessarily fit issues, you can use the forum
to connect with the community: https://gitlab.com/bloom42/forum


## Security

If you found a security issue affecting this project, please do not open a public issue and refer to our
[dedicated security page](https://bloom.sh/security) instead. Thank you!


## License

See `LICENSE.txt` and [https://bloom.sh/licensing](https://bloom.sh/licensing)

From an original work by [machinebox](https://github.com/machinebox): [graphql](https://github.com/machinebox/graphql) - commit 3a92531802258604bd12793465c2e28bc4b2fc85
