package models

type ListOptions struct {
	PageSize  int
	PageToken string
	UserId    string
	CourseId  string
}

func (opts *ListOptions) SetDefaults() {
	if opts.PageSize <= 0 {
		opts.PageSize = 10
	}
	if opts.PageToken == "" {
		opts.PageToken = "0"
	}
}
