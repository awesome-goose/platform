Declaration types:

- migration
- seed
- controller
- route
- middleware
- validator
- service

Module type:
type Module struct {
Imports []Module
Exports []Injectable
Declarations []Injectable
}

- Declaration, Exportable can be any of migration, seed, controller, route, middleware, validator or service
- provide dummy interface definition for each cmigration, seed, ontroller, route, middleware, validator, service

- on run
- - recursivelly go thru modules
- - - resolve and ingest migrations into container
- - - resolve and ingest seeds into container
- - - resolve and ingest models into container
- - - resolve and ingest controllers into container
- - - resolve and ingest routes into container
- - - resolve and ingest services into container
- - - resolve and ingest middlewares into container
- - - resolve and ingest middlewares into container
- - - resolve and ingest config into container
- return migration, seed, controller, route, middleware, validator or service

- reoslve rules
- - declarations in a module can only use values defined in their modules or values exported by imported modules

- extract migrations, resolve and run
- extract seeds, resolve and run
- extract routes and build route tree

- build route tree
- - extract route definition from all routes
- - merge into a single routes
- - transverse the routes and build a tree of expanded paths => handler
- - handler consist of ordered middleware, validator and controller handler
- - controller handler consist of a slice of [Controller Type, Method]
- - e.g for this path, call x method of y controller, pass a,b,c... to x. a,b,c are resolved from container base on type hint by user

- - during the transversal maintain a stack of OnBoot methods for any declarable that implements OnBoot
- - during the transversal maintain a stack of OnShutdown methods for any declarable that implements OnShutdown

- the Traverse method should return:
- - the container with all the resolved entities
- - the route tree with all the registered routes
- - each route should include all the middlewares, validators, and handlers
- - the stack of OnBoot methods
- - the stack of OnShutdown methods
