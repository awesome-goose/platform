package contracts

type Injectable interface{}

type Module struct {
	Imports      []Module
	Exports      []Injectable
	Declarations []Injectable
}
