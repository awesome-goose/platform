package contracts

type Response interface {
	Write(data []byte) error
}
