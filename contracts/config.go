package contracts

type Config interface {
	Dir() string
	Tree() map[string]any
	Get(path string) (any, error)
	Set(path string, value any) error
	Export(namespace string, config any) error
	Import(namespace string, in any) error
}
