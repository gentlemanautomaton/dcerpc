package ndr

type encTest1 struct {
	MaxLength int
	Length    int
	Data      []byte `idl:"size_is(MaxLength),length_is(Length)"`
}

type encTest2 struct {
	MaxLength int
	Length1   int
	Length2   int
	Data      [][]byte `idl:"size_is(MaxLength,20),length_is(Length1,Length2)"`
}
