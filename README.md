# kira

kira is a remote docker based sandbox code runner written in Go.

Currently it supports the following languages:

- Go (1.17)
- Java (8)
- C (Latest)
- Python (3.9.1)

## Installation

For the installation of kira, you need to have Docker and Go installed.

If you want to have the latest kira image on your machine, execute the [`build_kira_image.sh`](https://github.com/FlorianWoelki/kira/blob/main/build/build_kira_image.sh) script.

In addition, you need to pull the latest images by executing [`pull_images.sh`](https://github.com/FlorianWoelki/kira/blob/main/build/pull_images.sh). This will pull all the docker images that are being used by kira. This step should only be executed once.

## Usage

For now, you can only run your code by manipulating the [`main.go`](https://github.com/FlorianWoelki/kira/blob/main/main.go). You need to specify the language and the to be executed code.

This will create a new container sandbox and tries to execute the code you have passed.

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License

[MIT](https://choosealicense.com/licenses/mit/)