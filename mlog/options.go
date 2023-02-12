package mlog

type ConfigOptions struct {
	OutputPath    string
	MaxFileSizeMB int
	MaxBackups    int
	MaxAges       int
	Compress      bool
	LocalTime     bool
}

type Option func(*ConfigOptions)

func WithOutPutPath(outPutPath string) Option {
	return func(o *ConfigOptions) {
		o.OutputPath = outPutPath
	}
}

func WithMaxFileSizeMB(maxFileSizeMB int) Option {
	return func(o *ConfigOptions) {
		o.MaxFileSizeMB = maxFileSizeMB
	}
}

func WithMaxBackups(maxBackups int) Option {
	return func(o *ConfigOptions) {
		o.MaxBackups = maxBackups
	}
}

func WithMaxAges(maxAges int) Option {
	return func(o *ConfigOptions) {
		o.MaxAges = maxAges
	}
}

func WithCompress(compress bool) Option {
	return func(o *ConfigOptions) {
		o.Compress = compress
	}
}

func WithLocalTime(localTime bool) Option {
	return func(o *ConfigOptions) {
		o.LocalTime = localTime
	}
}
