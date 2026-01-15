package types

type Student struct {
	Id    int
	Name  string `validate:"required"` // Here now become this field required
	Email string `validate:"required"` // Here we can also use email, and some other validation error things
	Age   int    `validate:"required"` // These field validation use when create data using this Student struct
}
