package errors

const (
	Exists        = "already Exists"
	NotFound      = "does not Exist"
	DBError       = "connection Failed"
	AuthFailed    = "auth Failed"
	InvalidParams = "no valid parameters"

	OutOfRange = "out of range"
	TooLong    = "too long"
)

const (
	Id         = "id"
	Name       = "name"
	Preptime   = "preptime"
	Difficulty = "difficulty"
	Vegetarian = "vegetarian"
)

const (
	Rate   = "rate"
	RateId = "id"
)

const (
	Authentication = "Authentication"
	Bearer         = "bearer"
	RecipeId      = "id"
)

type ExistErr struct {
	s string
}
func (myErr *ExistErr) Error() string {
	return Exists
}

type NotFoundErr struct {
	s string
}
func (myErr *NotFoundErr) Error() string {
	return NotFound
}

type DBErr struct {
	s string

}
func (myErr *DBErr) Error() string {
	return DBError
}

type AuthFailedErr struct {
	s string

}
func (myErr *AuthFailedErr) Error() string {
	return AuthFailed
}

type InvalidParamsErr struct {
	s string

}
func (myErr *InvalidParamsErr) Error() string {
	return InvalidParams
}

