package ndr

type encTest1 struct {
	MaxLength int
	Length    int
	Data      []byte `idl:"size_is(MaxLength),length_is(Length)"`
}
