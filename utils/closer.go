package utils

import (
	"io"
)

func CloseSafe(closer io.Closer) {
	_ = closer.Close()
}
