package atexit

var hooks = []func(){}

func Register(f func()) {
	hooks = append(hooks, f)
}

func Exec() {
	for i := len(hooks)-1; i >= 0; i-- {
		hooks[i]()
	}
}
