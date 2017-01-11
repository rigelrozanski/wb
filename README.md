# whiteboard

_virtual whiteboard for terminal_

---

The intention behind this program is to provide a space for quickly retrieving and saving notes from terminal, I use it for saving obscure terminal commands

### Installation

1. Make sure you [have Go installed][1] and [put $GOPATH/bin in your $PATH][2]
2. [Install Cobra][3]
3. run `go get github.com/rigelrozanski/wb`
4. run `go install wb`

[1]: https://golang.org/doc/install
[2]: https://github.com/tendermint/tendermint/wiki/Setting-GOPATH 
[3]: https://github.com/spf13/cobra#installing

###  Usage

once installed run the command `wb` from any terminal window to read a common text file to terminal

### Command List
  
`wb --help` 	diplays program details and command list  
`wb` 		open the whiteboard	
`wb edit` 	edit the whiteboard	

### Contributing

1. Fork it
2. Create your feature branch (git checkout -b my-new-feature)
3. Commit your changes (git commit -am 'Add some feature')
4. Push to the branch (git push origin my-new-feature)
5. Create new Pull Request

### License

whiteboard is released under the Apache 2.0 license.
