package handlers

type Error struct {
	Msg string
}

var BadRequestError = Error{
	Msg: "bad request",
}
