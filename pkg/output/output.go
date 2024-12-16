package output

import (
	"fmt"
	"os"
)

func Error(message string) {
	fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}
