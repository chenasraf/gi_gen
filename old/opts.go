package main

type Behavior int

const (
	BehaviorOverwrite Behavior = iota
	BehaviorAppend
	BehaviorSkip
)

type Nullable[T any] struct {
	HasValue bool
	Value    T
}

func NewNullable[T any](hasValue bool, value T) Nullable[T] {
	return Nullable[T]{hasValue, value}
}

type GIGenOptions struct {
	Languages    *[]string
	CleanOutput  Nullable[bool]
	AutoDiscover Nullable[bool]
	Behavior     Nullable[Behavior]
}
