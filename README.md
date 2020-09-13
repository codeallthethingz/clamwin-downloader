# Clamwin DB downloader.

Downloads the clamwin DB, and extracts main.cvd from it.

Example usage

```golang
func main() {
	file, err := os.Create("virus.db")
	if err != nil {
		panic(err)
	}
	if err := NewClamwinConnector().Download(file); err != nil {
		panic(err)
	}
}
```