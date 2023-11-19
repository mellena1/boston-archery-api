package errors

// swagger:model Error
type Error struct {
	Msg string
}

var BadRequestError = Error{
	Msg: "bad request",
}

var UnauthorizedError = Error{
	Msg: "unauthorized",
}

var NotFoundError = Error{
	Msg: "not found",
}

var AlreadyExistsError = Error{
	Msg: "item already exists",
}
