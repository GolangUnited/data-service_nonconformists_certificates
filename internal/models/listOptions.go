package models

// ListOptions should be used as set of filters
// for any Database List method implementation
type ListOptions struct {
	Limit    int
	Offset   int
	UserId   string
	CourseId string
}

// SetDefaults checks and sets defaults
// for Limit and Offset
func (opts *ListOptions) SetDefaults() {
	if opts.Limit <= 0 {
		opts.Limit = 10
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}
}
