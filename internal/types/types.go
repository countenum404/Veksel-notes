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

type SpellingError struct {
	Code int      `json:"code"`
	Pos  int      `json:"pos"`
	Row  int      `json:"row"`
	Col  int      `json:"col"`
	Len  int      `json:"len"`
	Word string   `json:"word"`
	S    []string `json:"s"`
}

type SpellResult []SpellingError

type SpellResponse struct {
	NoteRequest CreateNoteRequest `json:"created_note"`
	Spells      SpellResult       `json:"spells"`
}
