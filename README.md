# whiteboard

_virtual whiteboard for terminal_

---

The intention behind this program is to provide a space for quickly retrieving
and saving notes from command line

### Installation

1. Make sure you [have Go installed][1] and [put $GOPATH/bin in your $PATH][2]
2. run `go get github.com/rigelrozanski/wb`
3. run `go install wb`

[1]: https://golang.org/doc/install
[2]: https://github.com/tendermint/tendermint/wiki/Setting-GOPATH 
[3]: https://github.com/spf13/cobra#installing

###  Example Usage

The following are a list of commands that can be run in terminal. wb can be run
while navigated to any directory. 

| Command    | Description                                     |
|------------|-------------------------------------------------|
| wb         | Opens the default whiteboard                    |
| wb cat     | Prints the default whiteboard                   |
| wb new foo | Create new wb named 'foo'                       |
| wb foo     | Open an existing whiteboard named 'foo'         |
| wb cat foo | Prints the contents of 'foo'                    |
| wb rm foo  | Deletes 'foo'                                   |
| wb ls      | List all whiteboards                            |
| wb log     | List all changes to whiteboards since last push |
| wb stats   | List repo stats, additions, deletions per wb    |

### Other Notes
 - a file in the root of this repo named config.txt can be used to setup a
   custom location for your whiteboards, but by default the text files are
   stored under the repo root folder. 
 - a file in the root of this repo named `push.sh` set's the custom commands to
   trigger when the command `wb push` is used. I use these to backup my wbs in
   a private git repository :)
 - shortcuts can be defined within the wb named `shortcuts`. Each shortcut 
   is defined on a new line as follows: `shortcut-name wb-name` 
 - The following are reserved words which can not be used for custom boards:
   - new
   - cat
   - rm
   - ls
   - log
   - stats

### License

whiteboard is released under the Apache 2.0 license.
