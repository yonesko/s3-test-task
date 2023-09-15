package model

import (
	"io"
)

type File struct {
	Name string
	Body io.Reader
}
