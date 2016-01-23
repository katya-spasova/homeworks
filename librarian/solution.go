package librarian

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"sync"
)

// Creates a new library
func NewLibrary(librarians int) Library {
	m := make(map[string]*Book)
	librariansChan := make(chan Librarian, librarians)
	return &FancyLibrary{books: m, librarians: librariansChan, mutex: sync.Mutex{}}
}

// Librarian working in the library
type Librarian struct {
	request  chan LibraryRequest
	response chan LibraryResponse
	library  *FancyLibrary
}

// The librarian waits for requests and helps borrowing a book,
// returning a book, getting information about a book
func (librarian *Librarian) serve() {
	go func() {
		for {
			message, ok := <-librarian.request
			if ok {
				isbn := message.GetISBN()
				requestType := message.GetType()
				switch requestType {
				case BorrowBook:
					libraryResponse := librarian.borrowBook(isbn)
					librarian.response <- &libraryResponse
				case ReturnBook:
					libraryResponse := librarian.returnBook(isbn)
					librarian.response <- &libraryResponse
				case GetAvailability:
					libraryResponse := librarian.getAvailability(isbn)
					librarian.response <- &libraryResponse
				default:
					librarian.response <- &CoolLibraryResponse{book: nil,
						err: errors.New("Невалидна заявка")}
				}
			} else {
				break
			}
		}
	}()
}

// Reduces the available count for book with given isbn
// Returns information about the book
// Returns error if the book is not in the library or all copies are taken
func (librarian *Librarian) borrowBook(isbn string) CoolLibraryResponse {
	librarian.library.mutex.Lock()
	defer librarian.library.mutex.Unlock()
	book, ok := librarian.library.books[isbn]
	var err error = nil
	if !ok {
		err = errors.New("Непозната книга " + isbn)
	} else {
		if book.availableCount == 0 {
			err = errors.New("Няма наличност на книга " + isbn)
		}

		book.availableCount -= 1
		librarian.library.books[isbn] = book
	}
	return CoolLibraryResponse{book: book, err: err}
}

// Increases the available count for book with given isbn
// Returns information about the book
// Returns error if the book is not in the library
// or all the books of this type are already in the library
func (librarian *Librarian) returnBook(isbn string) CoolLibraryResponse {
	librarian.library.mutex.Lock()
	book, ok := librarian.library.books[isbn]
	defer librarian.library.mutex.Unlock()
	var err error = nil
	if !ok {
		err = errors.New("Непозната книга " + isbn)
	} else {
		if book.availableCount == book.registeredCount {
			err = errors.New("Всички копия са налични " + isbn)
		}

		book.availableCount += 1
		librarian.library.books[isbn] = book
	}
	return CoolLibraryResponse{book: book, err: err}
}

// Returns information about the book
// Returns error if the book is not in the library
func (librarian *Librarian) getAvailability(isbn string) CoolLibraryResponse {
	librarian.library.mutex.Lock()
	defer librarian.library.mutex.Unlock()
	book, ok := librarian.library.books[isbn]
	var err error = nil
	if !ok {
		err = errors.New("Непозната книга " + isbn)
	}
	return CoolLibraryResponse{book: book, err: err}
}

type Person struct {
	FirstName string `xml:"first_name" json:"first_name"`
	LastName  string `xml:"last_name" json:"last_name"`
}

// The type for the books that are stored in the library
type Book struct {
	ISBN            string `xml:"isbn,attr"`
	Title           string `xml:"title"`
	Author          Person `xml:"author"`
	registeredCount int
	availableCount  int
}

// Returns a string representation of the book
func (book Book) String() string {
	return "[" + book.ISBN + "] " + book.Title + " от " + book.Author.FirstName +
		" " + book.Author.LastName
}

// The library
type FancyLibrary struct {
	books      map[string]*Book
	librarians chan Librarian
	mutex      sync.Mutex
}

// Adds a book to the library
// Returns the count of all available copies in the library
// Return an error of the number of copies is more than 4
func (library *FancyLibrary) addBook(book *Book) (int, error) {
	// Assuming that books with same ISBN are the same
	library.mutex.Lock()
	defer library.mutex.Unlock()
	same_book := library.books[book.ISBN]
	if same_book != nil {
		// there are already 4 copies of the book - return error
		if same_book.registeredCount == 4 {
			return 4, errors.New("Има 4 копия на книга " + book.ISBN)
		}
		same_book.registeredCount += 1
		same_book.availableCount += 1
	} else {
		book.registeredCount = 1
		book.availableCount = 1
		library.books[book.ISBN] = book
	}
	return library.books[book.ISBN].availableCount, nil
}

// Add a book from json
// Returns the count of all available copies in the library
// Return an error of the number of copies is more than 4
func (library *FancyLibrary) AddBookJSON(data []byte) (int, error) {
	book := new(Book)
	err := json.Unmarshal(data, book)
	if err != nil {
		return 0, err
	}
	return library.addBook(book)
}

// Add a book from xml
// Returns the count of all available copies in the library
// Return an error of the number of copies is more than 4
func (library *FancyLibrary) AddBookXML(data []byte) (int, error) {
	book := new(Book)
	err := xml.Unmarshal(data, book)
	if err != nil {
		return 0, err
	}
	return library.addBook(book)
}

// Gets a free librarian to serve requests
// Librarians are fixed - passed as an argument to NewLibrary func
// Blocks if all librarians are busy
// Returns two channels
// write channel -  for sending of requests
// read channel - for receiving results
// On closing the request (write channel) the librarian is released
func (library *FancyLibrary) Hello() (chan<- LibraryRequest, <-chan LibraryResponse) {
	// The librarians can get up to 100 requests before someone reads the response
	request := make(chan LibraryRequest, 100)
	response := make(chan LibraryResponse, 100)
	librarian := Librarian{request: request, response: response, library: library}
	library.librarians <- librarian
	librarian.serve()
	return librarian.request, librarian.response
}

// Available types of requests
const (
	_ = iota
	BorrowBook
	ReturnBook
	GetAvailability
)

// The type for the requests the librarians work with
type CoolLibraryRequest struct {
	isbn        string
	requestType int
}

// Return the type of the request:
// 1 - Borrow book
// 2 - Return book
// 3 - Get availability information about book
func (request *CoolLibraryRequest) GetType() int {
	return request.requestType
}

// Return the isbn of the book for which is the request
func (request *CoolLibraryRequest) GetISBN() string {
	return request.isbn
}

// The type for the requests the librarians work with
type CoolLibraryResponse struct {
	book *Book
	err  error
}

// gets book content, an object implementing Stringer
// If the book does not exist the first result is nil
// returns an error if it has occured
// when "Return book", the content is not attached
func (response *CoolLibraryResponse) GetBook() (fmt.Stringer, error) {
	if response.err != nil {
		return nil, response.err
	}
	if response != nil && response.book != nil {
		return response.book, nil
	} else {
		return nil, errors.New("Празен отговор")
	}
}

// available - how many books are available after the request is performed
// registered - how many copies are registered from this book (маx 4).
func (response *CoolLibraryResponse) GetAvailability() (available int, registered int) {
	if response != nil && response.book != nil {
		return response.book.availableCount, response.book.registeredCount
	} else {
		return 0, 0
	}
}
