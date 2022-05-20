# kira

kira is a remote docker based sandbox code execution engine written in Go.

Currently it supports the following languages:

- Go (golang:1.17-alpine)
- Java (openjdk:8u232-jdk)
- C (gcc:latest)
- C++ (gcc:latest)
- Python (python:3.9.1-alpine)
- JavaScript (node:lts-alpine)

## Installation

For the installation of kira, you need to have Docker and Go installed.

It is required for the sandbox environment to have a all-in-one image that includes all the functionality for executing the code in a docker container. For that, you need to create and build the image in `build/all-in-one-ubuntu`. You can do that by executing the `create-image` command in the `Makefile` or by executing:

```sh
docker build build/all-in-one-ubuntu -t all-in-one-ubuntu
```

## Usage

You can feel free to run the CLI by executing the `main.go` file in the root directory with the following command:

```sh
$ go run main.go
```

This will prompt you with some flags and commands you can use.

In addition, you can start the REST API running the `main.go` file in the `rest` directory. This will start the REST API on port `9090`.

### Commands and Flags for the CLI

The following section contains all the commands and flags that can be used while running the CLI.

<details>
  <summary>execute</summary>

  <p>
    The execute command will execute code in a containerized sandbox.
  </p>

  | Flag | Aliases | Description | Default |
  |---|---|---|---|
  | --language | -l, -lang | Set the language for the kira sandbox runner. | python |
  | --main | -m | Set the main file that should be executed first. | example code in runner struct |
  | --dir | -d | Set the specific directory that should be executed. | example code in runner struct |
</details>

### REST API endpoints

The following section contains all the REST API endpoints. The JSON body and endpoints follow the CLI structure.

<details>
  <summary>/execute</summary>

  <p>
    The execute endpoint will execute code in a containerized sandbox.
  </p>

  This JSON structure is an example for the request body:
  ```json
  {
      "language": "python",
      "content": "print(\"42 Hello World\")"
  }
  ```
</details>

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)