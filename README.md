# Clamwin DB downloader

Downloads the clamwin DB, and extracts main.cvd from it.

## Usage

```golang
package main

import (
	"os"

	"github.com/codeallthethingz/clamwin-downloader/clamwin"
)

func main() {
	file, err := os.Create("virus.db")
	if err != nil {
		panic(err)
	}
	if err := clamwin.NewClamwinConnector().Download(file); err != nil {
		panic(err)
	}
}
```


## License

Please note, clamwin DB is GNU GENERAL PUBLIC LICENSE, Version 3, and thus you cannot use this code in any commercial software.
