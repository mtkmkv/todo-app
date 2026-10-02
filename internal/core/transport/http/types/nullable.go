package core_http_types

import (
	"bytes"
	"encoding/json"

	"github.com/mtkmkv/todo-app/internal/core/domain"
)

type Nullable[T any] struct {
	domain.Nullable[T]
}

// To Domain преобразует транспортный Nullable в доменный domain.Nullable
func (n Nullable[T]) ToDomain() domain.Nullable[T] {
	return domain.Nullable[T]{
		Value: n.Value,
		Set:   n.Set,
	}
}

// UnmarshalJSON распаковывает JSON с отслеживанием явной передачи поля и значения null
func (n *Nullable[T]) UnmarshalJSON(b []byte) error {
	n.Set = true

	if bytes.Equal(b, []byte("null")) {
		n.Value = nil
		return nil
	}

	var value T
	if err := json.Unmarshal(b, &value); err != nil {
		return err
	}

	n.Value = &value
	return nil
}

// MarshalJSON сериализует значение обратно в компактный JSON
func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.Set || n.Value == nil {
		return []byte("null"), nil
	}
	return json.Marshal(n.Value)
}