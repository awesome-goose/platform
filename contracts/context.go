package contracts

type Context interface {
	Request() Request
	Response() Response

	Segments() []string
}
