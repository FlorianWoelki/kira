# kira

kira is a remote docker based sandbox code execution engine written in Go.

Currently it supports the following languages:

- Go (golang:1.17-alpine)
- Java (openjdk:8u232-jdk)
- C (gcc:latest)
- C++ (gcc:latest)
- Python (python:3.9.1-alpine)
- JavaScript (node:lts-alpine)
- JavaScript (node:lts-alpine & latest tsc)
- Julia (julia:1.7.1-alpine)
- Elixir (elixir:1.13.1-alpine)
- Swift (swift:5.5.2)

## Installation

For the installation of kira, you need to have Docker and Go installed.

If you want to have the latest kira image on your machine, execute the [`build_kira_image.sh`](https://github.com/FlorianWoelki/kira/blob/main/build/build_kira_image.sh) script.

In addition, you need to pull the latest images by executing [`pull_images.sh`](https://github.com/FlorianWoelki/kira/blob/main/build/pull_images.sh). This will pull all the docker images that are being used by kira. This step should only be executed once.

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
  <summary>Execute</summary>

  <p>
    The execute command will execute code in a containerized sandbox.
  </p>

  | Flag | Aliases | Description | Default |
  |---|---|---|---|
  | --language | -l, -lang | Set the language for the kira sandbox runner. | python |
  | --file | -f | Set the specific file that should be executed. | example code in runner struct |
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