package temple

import "errors"

// errRefused is the error every refusing function here gives.
var errRefused = errors.New("refused")

type feintThing struct{ id string }

func emptyLiteral() (feintThing, error) {
	return feintThing{}, errRefused
}

func emptyDeclared() (feintThing, error) {
	var made feintThing
	return made, errRefused
}

func fieldsSet(id string) (feintThing, error) {
	return feintThing{id: id}, errRefused
}

func vowedEmpty() (feintThing, error) {
	//ninja:vow feint the dojo keeps one vowed return to show it
	return feintThing{}, errRefused
}

func basicOnly() (int, error) {
	return 0, errRefused
}

func nilPointer() (*feintThing, error) {
	return nil, errRefused
}

func nilSlice() ([]feintThing, error) {
	return nil, errRefused
}

func noRefusal() (feintThing, error) {
	return feintThing{}, nil
}

func setThenReturned(id string) (feintThing, error) {
	var made feintThing
	made.id = id
	return made, errRefused
}
