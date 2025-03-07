package model

type WebResponse[T any] struct {
	Data   T     `json:"data"`
	Errors error `json:"errors,omitempty"`
	Status int
}
