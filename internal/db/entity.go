package db

type BaseModel[T any] interface {
	ConvertToModel() T
}
