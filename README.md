Tool for the automatic creation of constructors using the functional options-pattern

## What is opt?

The Opt tool generates constructors from struct types. It saves the user the cration of boiler plate code when creating With-functions.

## Installation

Install opt using the go install command:

```bash
go install github.com/boecklim/opt
```

## How to use
The functional options pattern is very useful as it allows for flexible constructors where not each usage of a constructor needs to give all parameters. Here is a good article which describes the functional options pattern: https://uptrace.dev/blog/golang-functional-options.html

In order to generate a constructor for a type, add a go-generate magic comment (or use `opt` as a cli tool).

`opt` can create a constructor with both optional and required parameters. A tag `opt:"true"` added after one struct member will indicate opt that for each members after that and including that member a `With...` function will be created. All members before that tag will be generated as required parameters. The output file will be in the same package and folder.

Example:
```
//go:generate opt -out constructor.go -rm . ExampleStruct
type ExampleStruct struct {
	firstMember  string
	SecondMember int
	timestamp    time.Time

	byteReader func() io.ByteReader `opt:"true"`
	aSlice     []float64
	aPointer   *string
}
```

This will create a file named `constructor.go` with the following content:
```
package example

import (
	"io"
	"time"
)

type Option func(i *ExampleStruct)

// With byteReader of type func() io.ByteReader
func WithByteReader(byteReader func() io.ByteReader) Option {
	return func(s *ExampleStruct) {
		s.byteReader = byteReader
	}
}

// With aSlice of type []float64
func WithASlice(aSlice []float64) Option {
	return func(s *ExampleStruct) {
		s.aSlice = aSlice
	}
}

// With aPointer of type *string
func WithAPointer(aPointer *string) Option {
	return func(s *ExampleStruct) {
		s.aPointer = aPointer
	}
}

func New(firstMember string, SecondMember int, timestamp time.Time, opts ...Option) *ExampleStruct {
	newInstance := ExampleStruct{}

	for _, opt := range opts {
		opt(&newInstance)
	}

	return &newInstance
}
```

Using the `-rm` flag will remove a generated file of the same if it already exists.

## Disclaimer
This tool is still in development and may contain bugs.

## Acknowledgements
Special thanks to Mat Ryer and David Hernandez. This module has been much inspired by [moq](https://github.com/matryer/moq)
