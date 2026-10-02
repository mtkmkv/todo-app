package domain

/*
	Помогает отличать пустой JSON от переданного null в JSON
*/
type Nullable[T any] struct {
	Value *T
	Set   bool
}
