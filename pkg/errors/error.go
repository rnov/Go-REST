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
	RecipeId       = "id"
)

type ExistErr struct {
	msg string
}

func (myErr *ExistErr) Error() string {
	return Exists
}

func NewExistErr(message string) *ExistErr {
	return &ExistErr{
		msg: message,
	}
}

type NotFoundErr struct {
	msg string
}

func (myErr *NotFoundErr) Error() string {
	return NotFound
}

func NewNotFoundErr(message string) *NotFoundErr {
	return &NotFoundErr{
		msg: message,
	}
}

type DBErr struct {
	msg string
}

func (myErr *DBErr) Error() string {
	return DBError
}

func NewDBErr(message string) *DBErr {
	return &DBErr{
		msg: message,
	}
}

type AuthFailedErr struct {
	msg string
}

func (myErr *AuthFailedErr) Error() string {
	return AuthFailed
}

func NewAuthFailedErr(message string) *AuthFailedErr {
	return &AuthFailedErr{
		msg: message,
	}
}

type InvalidParamsErr struct {
	msg string
	Prm map[string]string
}

func (myErr *InvalidParamsErr) Error() string {
	return InvalidParams
}

func NewInvalidParamsErr(params map[string]string) *InvalidParamsErr {
	return &InvalidParamsErr{
		Prm: params,
	}
}


