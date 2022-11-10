package models

type ListOptions struct {
	Limit    int
	Offset   int
	UserId   string
	CourseId string
}

func (opts *ListOptions) SetDefaults() {
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}
}
