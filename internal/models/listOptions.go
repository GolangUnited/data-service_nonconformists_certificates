package models

// ListOptions should be used as set of filters
// for any Database List method implementation
type ListOptions struct {
	// Limits output to specified amount of entries
	Limit int

	// Offset in DB
	Offset int

	// Filter by UserId
	UserId string

	// Filter by CourseId
	CourseId string

	// Filter out deleted records
	ShowDeleted bool
}

// SetDefaults checks and sets defaults
// for Limit and Offset
func (opts *ListOptions) SetDefaults() {
	switch {
	// Default page size is 10
	case opts.Limit <= 0:
		opts.Limit = 10
	// Max page size is 100
	case opts.Limit > 100:
		opts.Limit = 100
	}
	if opts.Offset < 0 {
		opts.Offset = 0
	}
}
