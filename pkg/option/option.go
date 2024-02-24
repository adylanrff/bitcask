package option

const (
	DefaultSegmentSizeLimit uint = 1 << 10
)

type Options struct {
	SegmentSizeLimit uint
}

func (o *Options) Apply(opts ...OptionFunc) {
	for _, opt := range opts {
		opt(o)
	}
}

func NewOptions() *Options {
	return &Options{}
}

func NewDefaultOptions() *Options {
	return &Options{
		SegmentSizeLimit: DefaultSegmentSizeLimit,
	}
}

type OptionFunc func(*Options)
