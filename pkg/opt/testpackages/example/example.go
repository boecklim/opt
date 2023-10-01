package example

import (
	"io"
	"time"
)

//go:generate opt -out example_mock.go -rm . ExampleStruct
type ExampleStruct struct {
	firstMember  string
	SecondMember int
	timestamp    time.Time

	byteReader func() io.ByteReader `opt:"true"`
	aSlice     []float64
	aPointer   *string
	// AnInterface  http.CloseNotifier --> error
}
