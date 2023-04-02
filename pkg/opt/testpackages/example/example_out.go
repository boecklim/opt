// Code generated by opt; DO NOT EDIT.
// github.com/boecklim/opt

package example

// TODO: imports

type Option func(i *ExampleStruct)

// With firstMember of type string
func WithfirstMember(firstMember string) Option {
	return func(s *ExampleStruct) {
		s.firstMember = firstMember
	}
}

// With SecondMember of type int
func WithSecondMember(SecondMember int) Option {
	return func(s *ExampleStruct) {
		s.SecondMember = SecondMember
	}
}

func New(opts ...Option) *ExampleStruct {
	newInstance := ExampleStruct{}

	for _, opt := range opts {
		opt(&newInstance)
	}

	return &newInstance
}