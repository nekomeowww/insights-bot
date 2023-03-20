package options

type CallOptions[O any] struct {
	applyFunc func(o *O)
}

func NewCallOptions[O any](applyFunc func(o *O)) CallOptions[O] {
	return CallOptions[O]{
		applyFunc: applyFunc,
	}
}

func ApplyCallOptions[O any](callOptions []CallOptions[O], defaultOptions ...O) *O {
	o := new(O)
	if len(defaultOptions) > 0 {
		*o = defaultOptions[0]
	}

	for _, callOption := range callOptions {
		callOption.applyFunc(o)
	}

	return o
}
