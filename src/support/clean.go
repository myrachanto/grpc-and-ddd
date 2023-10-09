package support

import "os"

func Cleaner(filename string) {
	if filename != "" && filename != "undefined" {
		os.Remove("./src/public" + filename)
	}
}
