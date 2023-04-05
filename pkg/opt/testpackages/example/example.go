package example

import (
	"io"
	"time"
)

type ExampleStruct struct {
	firstMember  string
	SecondMember int
	timestamp    time.Time
	byteReader   func() io.ByteReader
	aSlice       []float64
	aPointer     *string
	// AnInterface  http.CloseNotifier --> error
}
