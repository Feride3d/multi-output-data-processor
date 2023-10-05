// Entity represents models of the application.
package entity

// InputData represents the input data for processing.
type InputData struct {
	Tag  string `json:"tag" validate:"required"`
	Data string `json:"data" validate:"required"`
}

// OutputChannel represents an output channel.
type OutputChannel struct {
	Name interface{}
}
