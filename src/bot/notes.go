package main


type Note struct {
	title string
	body string
	class string
	tags []string
}

func add_note(note Note) {
	// store somewhere
}

func read_notes(note Note) {
	// store somewhere
}

func find_notes_by_class(note Note) {
	// store somewhere
}

func add_secure_note(note Note) {
	// store somewhere
}

type Selector struct {
	tags []string
}

func read_secure_notes(selector Selector) {
	// store somewhere
}
