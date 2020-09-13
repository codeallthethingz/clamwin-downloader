package main

import (
	"os"

	"github.com/codeallthethingz/clamwin-downloader/clamwin"
)

// entry point if you're not using this as a package
func main() {
	file, err := os.Create("virus.db")
	if err != nil {
		panic(err)
	}
	if err := clamwin.NewClamwinConnector().Download(file); err != nil {
		panic(err)
	}
}
