# File Store
keep files on a remote server and operate via local client

## Download
### Container images (Server)
`master` branch images are tagged as `:master`<br>
Both released and master branch container images are published in to the [quay.io registry](https://quay.io/repository/jkandasa/file-store-server)
#### To Run Container Image
```bash
podman run --detach --name store-server \
    --volume $PWD/_store:/app/_store \
    --publish 8080:8080 \
    --env  TZ="Asia/Kolkata" \
    --restart unless-stopped \
    quay.io/jkandasa/file-store-server:master
```

### Download Executables
* [Released versions](https://github.com/jkandasa/file-store/releases)
* [Pre Release](https://github.com/jkandasa/file-store/releases/tag/master) - `master` branch executables


## Configuration
### Server
Server can be started with container image or with an executable binary. 
* For container images refer [Container Image](#to-run-container-image)
* Binary images can be downloaded from [release page](https://github.com/jkandasa/file-store/releases)
  * To execute binary download and extract `store-server-*` for your platform.

```bash
# to get server version details
$ ./store-server -version

version:master, buildDate:2022-11-21T01:54:41+00:00, gitCommit:d8bd781d2096cd8f8565de812ced6d5109df177c, goLang:go1.19, platform:linux/amd64}

# to execute
$ ./store-server -port 8080

2022-11-21T08:02:32.642+0530	info	server/main.go:36	version details	{"version": "{version:master, buildDate:2022-11-21T01:54:41+00:00, gitCommit:d8bd781d2096cd8f8565de812ced6d5109df177c, goLang:go1.19, platform:linux/amd64}"}
2022-11-21T08:02:32.643+0530	info	server/main.go:46	listening HTTP service on	{"address": "0.0.0.0:8080"}

```
### Client
Client is available only in binary executable format.
To execute binary download and extract `store-client-*` for your platform.
* To change server address update `STORE_SERVER` environment variable. default: `http://127.0.0.1:8080`

```bash
# update server address
export STORE_SERVER=http://127.0.0.1:8080

$ ./store version
```

To get command helps in the client, always use `-h` or `--help` in the end of command, example:
```bash
$ ./store --help
Storage Client
  
This client helps you to control your storage server from the command line.

Usage:
  storage [command]

Available Commands:
  add         Adds the files to remote server
  completion  Generate the autocompletion script for the specified shell
  freq-words  Prints the frequent words count from the available text files
  help        Help about any command
  ls          Prints the available file details from the server
  rm          Removes the files from the remote server
  update      Updates/sync the files on the remote server
  version     Print the client version information
  wc          Prints the word count from the available text files on the remote server

Flags:
  -h, --help            help for storage
      --hide-header     hides the header on the console output
      --insecure        connect to server in insecure mode (default true)
  -o, --output string   output format. options: yaml, json, console (default "console")
      --pretty          JSON pretty print

Use "storage [command] --help" for more information about a command.


$ ./store add --help
Adds the files to remote server

Usage:
  storage add [flags]

Examples:
  # add files
  store add file1.txt file2.txt

Flags:
  -h, --help   help for add

Global Flags:
      --hide-header     hides the header on the console output
      --insecure        connect to server in insecure mode (default true)
  -o, --output string   output format. options: yaml, json, console (default "console")
      --pretty          JSON pretty print

```