package errors

type Error struct {
	Msg string
}

var BadRequestError = Error{
	Msg: "bad request",
}

var UnauthorizedError = Error{
	Msg: "unauthorized",
}
