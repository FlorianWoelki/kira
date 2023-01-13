# kira

kira is a remote docker based sandbox code execution engine written in Go.

Currently it supports the following languages:

`python`, `javascript`, `golang`, `java`, `bash`, `c`, `c++`

## Installation

For the installation of kira, you need to have Docker and Go installed.

It is required for the sandbox environment to have a docker image that exposes an API that includes all the functionality for executing the code. For that, you need to execute the following command to build the image and start the container:

```sh
docker-compose up
```

## Usage

You can feel free to run the REST API by executing the `main.go` file in the `rest` directory with the following command:

```sh
$ go run rest/main.go
```

The REST API will start on port `:9090`.

### REST API endpoints

The following section contains all the REST API endpoints. The JSON body and endpoints follow the CLI structure.

<details>
  <summary>POST /execute</summary>

  <p>
    The execute endpoint will execute code in a containerized sandbox. Tests for the
    printed output can also be specified.
  </p>

  This JSON structure is an example for the request body:
  ```json
  {
      "language": "python",
      "content": "print(\"42 Hello World\")",
      "tests" [
        { "name": "First test case ", "actual": "42 Hello World" },
        { "name": "First test case ", "actual": "41 Hello World" }
      ]
  }
  ```

  You can also add an optional query parameter called `bypass_cache` and set it to `true`,
  if you want to bypass the cache.
</details>

<details>
  <summary>GET /languages</summary>

  <p>
    Will return all languages that are possible for remote execution.
  </p>

  This JSON structure is an example for the response body:
  ```json
  [
      {
          "name": "python",
          "version": "3.7.10",
          "extension": ".py",
          "timeout": 10
      },
      {
          "name": "javascript",
          "version": "16.3.1",
          "extension": ".js",
          "timeout": 10
      },
      // ...
  ]
  ```
</details>

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)