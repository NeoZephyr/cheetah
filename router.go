package main

import "cheetah/framework"

func registerRouter(core *framework.Core) {
	core.Get("foo", FooController)
}
