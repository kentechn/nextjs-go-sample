package todo

// ListInput selects which todos to return.
type ListInput struct {
	// Done filters by completion state; nil means "any".
	Done *bool
}

// CreateInput is the raw input of the create use case.
type CreateInput struct {
	Title string
}
