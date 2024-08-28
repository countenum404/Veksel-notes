package types

type ApiError struct {
	Error string
}

type User struct {
	ID        int64
	Fisrtname string
	Lastname  string
	Username  string
	Password  string
}

type Note struct {
	ID      int64  `json:"id"`
	Header  string `json:"header"`
	Content string `json:"content"`
}

type CreateNoteRequest struct {
	Header  string `json:"header"`
	Content string `json:"content"`
}
