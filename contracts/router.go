package contracts

type Route struct {
	Path    string // even segments like method should be included in the path
	Handler func(Context) any

	Validators  Validators
	Middlewares Middlewares

	Children []Route

	Params map[string]string
}

type Routes = []Route

type Router interface {
	Find(routes Routes, segments []string) (*Route, error)
}
