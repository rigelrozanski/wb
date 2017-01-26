# whiteboard

_virtual whiteboard for terminal_

---

The intention behind this program is to provide a space for quickly retrieving and saving notes from terminal, I use it for saving obscure terminal commands

### Installation

1. Make sure you [have Go installed][1] and [put $GOPATH/bin in your $PATH][2]
2. run `go get github.com/rigelrozanski/wb`
3. run `go install wb`

[1]: https://golang.org/doc/install
[2]: https://github.com/tendermint/tendermint/wiki/Setting-GOPATH 
[3]: https://github.com/spf13/cobra#installing

###  Example Usage

The following are a list of commands that can be run in terminal. wb can be run while navigated to any directory. 

| Command   | Alt. Cmd. | Description                               |
|-----------|-----------|-------------------------------------------|
| wb        |           | Opens the default whiteboard              |
| wb ed     |           | Edit the default whiteboard               |
| wb foo    |           | View an existing whiteboard named 'foo'   |
| wb nu foo | wb foo nu | Create new wb named 'foo'                 |
| wb ed foo | wb foo ed | Edit an existing whiteboard named 'foo'   |
| wb rm foo | wb foo rm | Remove an existing whiteboard named 'foo' |
| wb list   |           | List all whiteboards                      |

### Other Notes

 - Raw text files are stored under the repo root folder
 - The following are reserved words which can not be used for custom boards:
   - nu
   - ed
   - rm
   - list

### Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

### License

whiteboard is released under the Apache 2.0 license.
