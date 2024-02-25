package option

import (
	"encoding/json"
	"fmt"
)

// Option values explicitly indicate the presence or absence of a value.
// This generic type can be used as replacement for nil values.
type Option[T any] struct {
	value   T
	present bool
}

// Some creates a present Option.
func Some[T any](v T) Option[T] {
	return Option[T]{value: v, present: true}
}

// None creates an absent Option.
func None[T any]() Option[T] {
	return Option[T]{value: zeroValueOf[T](), present: false}
}

// String returns a string representation of the Option.
func (opt Option[T]) String() string {
	if opt.Absent() {
		return "None"
	}
	return fmt.Sprintf("Some(%#v)", opt.value)
}

// Value unwraps the value of the Option.
// It panics if the Option is absent. Otherwise, it returns the value.
func (opt Option[T]) Value() T {
	if !opt.Present() {
		panic("option: called Value() on an absent option")
	}
	return opt.value
}

// Absent returns true if the Option is absent.
func (opt Option[T]) Absent() bool { return !opt.present }

// Present returns true if the Option is present.
func (opt Option[T]) Present() bool { return opt.present }

// ValueOr returns the value of the Option if present, otherwise, it returns the fallback value.
func (opt Option[T]) ValueOr(fallback T) T {
	if opt.Absent() {
		return fallback
	}
	return opt.value
}

// ValueOrBy returns the value of the Option if present, otherwise, it returns the value from the supplier.
func (opt Option[T]) ValueOrBy(supplier func() T) T {
	if opt.Absent() {
		return supplier()
	}
	return opt.value
}

// Alt returns the Option if present, otherwise, it returns the alternative Option.
func (opt Option[T]) Alt(alt Option[T]) Option[T] {
	if opt.Absent() {
		return alt
	}
	return opt
}

// AltBy returns the Option if present, otherwise, it returns the Option from the supplier.
func (opt Option[T]) AltBy(supplier func() Option[T]) Option[T] {
	if opt.Absent() {
		return supplier()
	}
	return opt
}

func (opt Option[T]) MarshalJSON() ([]byte, error) {
	if opt.Absent() {
		return json.Marshal(jsonValue[T]{Kind: "Option::None"})
	}
	return json.Marshal(jsonValue[T]{Kind: "Option::Some", Value: opt.value})
}

func (opt *Option[T]) UnmarshalJSON(data []byte) error {
	var jv jsonValue[T]
	if err := json.Unmarshal(data, &jv); err != nil {
		return err
	}

	switch jv.Kind {
	case "Option::None":
		*opt = None[T]()
	case "Option::Some":
		*opt = Some(jv.Value)
	default:
		return fmt.Errorf("option: unexpected kind: %q", jv.Kind)
	}
	return nil
}

func zeroValueOf[T any]() (z T) { return }

type jsonValue[T any] struct {
	Kind  string `json:"kind"` // "Some" or "None"
	Value T      `json:"value"`
}
