// Code generated by opt; DO NOT EDIT.
// github.com/boecklim/opt

package example

type Option func(i *ExampleStruct)

func New(opts ...Option) *ExampleStruct {
	newInstance := ExampleStruct{}

	for _, opt := range opts {
		opt(&newInstance)
	}

	return &newInstance
}
