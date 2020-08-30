package pkg

type Author struct {
	Name string
}

type Book struct {
	Author Author
	Title  string
}

type Library struct {
	Books []Book
}
