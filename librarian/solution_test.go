package librarian

import (
	"errors"
	"testing"
)

func TestGetBook(t *testing.T) {
	author := Person{FirstName: "Anno", LastName: "Birkin"}
	book := Book{ISBN: "0954540018", Author: author,
		registeredCount: 2, availableCount: 1,
		Title: "Who Said the Race is Over?"}
	response := CoolLibraryResponse{book: &book, err: nil}
	found, _ := response.GetBook()
	expected := "[0954540018] Who Said the Race is Over? от Anno Birkin"
	if found.String() != expected {
		t.Errorf("Expected\n---\n%s\n---\nbut found\n---\n%s\n---\n", expected, found.String())
	}
}

func TestGetBookNoBook(t *testing.T) {
	response := CoolLibraryResponse{}
	found, err := response.GetBook()
	expectedError := "Празен отговор"
	if found != nil {
		t.Errorf("Expected\n---\nnil\n---\nbut found\n---\n%s\n---\n", found)
	}
	if err.Error() != expectedError {
		t.Errorf("Expected\n---\n%s\n---\nbut found\n---\n%s\n---\n", expectedError, err.Error())
	}
}

func TestGetBookError(t *testing.T) {
	response := CoolLibraryResponse{err: errors.New("Просто грешка")}
	found, err := response.GetBook()
	expectedError := "Просто грешка"
	if found != nil {
		t.Errorf("Expected\n---\nnil\n---\nbut found\n---\n%s\n---\n", found)
	}
	if err.Error() != expectedError {
		t.Errorf("Expected\n---\n%s\n---\nbut found\n---\n%s\n---\n", expectedError, err.Error())
	}
}

// available - how many books are available after the request is performed
// registered - how many copies are registered from this book (маx 4).
func TestGetAvailability(t *testing.T) {
	author := Person{FirstName: "Anno", LastName: "Birkin"}
	book := Book{ISBN: "0954540018", Author: author,
		registeredCount: 2, availableCount: 1,
		Title: "Who Said the Race is Over?"}
	response := CoolLibraryResponse{&book, nil}
	available, registered := response.GetAvailability()
	availableExpected := 1
	registeredExpected := 2
	if available != availableExpected {
		t.Errorf("Expected available\n---\n%s\n---\nbut found\n---\n%s\n---\n", availableExpected, available)
	}
	if registered != registeredExpected {
		t.Errorf("Expected registered\n---\n%s\n---\nbut found\n---\n%s\n---\n", registeredExpected, registered)
	}
}

// available - how many books are available after the request is performed
// registered - how many copies are registered from this book (маx 4).
func TestNoBookGetAvailability(t *testing.T) {
	response := CoolLibraryResponse{}
	available, registered := response.GetAvailability()
	availableExpected := 0
	registeredExpected := 0
	if available != availableExpected {
		t.Errorf("Expected available\n---\n%s\n---\nbut found\n---\n%s\n---\n", availableExpected, available)
	}
	if registered != registeredExpected {
		t.Errorf("Expected registered\n---\n%s\n---\nbut found\n---\n%s\n---\n", registeredExpected, registered)
	}
}

func TestGetTypeBorrow(t *testing.T) {
	request := CoolLibraryRequest{isbn: "8798797989vvg", requestType: BorrowBook}
	expected := 1
	found := request.GetType()
	if found != expected {
		t.Errorf("Expected available\n---\n%s\n---\nbut found\n---\n%s\n---\n", expected, found)
	}
}

func TestGetTypeReturn(t *testing.T) {
	request := CoolLibraryRequest{isbn: "8798797989vvg", requestType: ReturnBook}
	expected := 2
	found := request.GetType()
	if found != expected {
		t.Errorf("Expected available\n---\n%s\n---\nbut found\n---\n%s\n---\n", expected, found)
	}
}

func TestGetTypeGetAvailability(t *testing.T) {
	request := CoolLibraryRequest{isbn: "8798797989vvg", requestType: GetAvailability}
	expected := 3
	found := request.GetType()
	if found != expected {
		t.Errorf("Expected type\n---\n%s\n---\nbut found\n---\n%s\n---\n", expected, found)
	}
}

func TestGetISBN(t *testing.T) {
	request := CoolLibraryRequest{isbn: "8798797989vvg", requestType: GetAvailability}
	expected := "8798797989vvg"
	found := request.GetISBN()
	if found != expected {
		t.Errorf("Expected isbn\n---\n%s\n---\nbut found\n---\n%s\n---\n", expected, found)
	}
}

func TestAddBookJson(t *testing.T) {
	library := NewLibrary(5)

	str := `{
	  "isbn": "9781617293092",
	  "title": "Learn Go",
	  "author": {
		"first_name": "Nathan",
		"last_name": "Youngman"
	  },
	  "ratings": [5, 4, 4, 5, 1]
	}`
	count, err := library.AddBookJSON([]byte(str))
	if err != nil {
		t.Fatalf("An error occured while parsing json %s", err.Error())
	}

	if count != 1 {
		t.Errorf("Count must be 1 but is %d", count)
	}

	switch l := library.(type) {
	case *FancyLibrary:
		b := l.books["9781617293092"]
		if b == nil {
			t.Fatalf("Book is nil, but it is expected to be in library")
		}

		if b.registeredCount != 1 {
			t.Errorf("Registered count of book should be 1, but is %s", b.registeredCount)
		}

		if b.availableCount != 1 {
			t.Errorf("Registered count of book should be 1, but is %s", b.availableCount)
		}

		if b.Author.FirstName != "Nathan" {
			t.Errorf("Author's first name must be Nathan, but is %s", b.Author.FirstName)
		}

		if b.Author.LastName != "Youngman" {
			t.Errorf("Author's last name must be Youngman, but is %s", b.Author.LastName)
		}

		if b.Title != "Learn Go" {
			t.Errorf("Title must be Learn Go, but is %s", b.Title)
		}
	}
}

func TestAddBookXml(t *testing.T) {
	library := NewLibrary(5)

	str := `
		<book isbn="0954540018">
			<title>Who said the race is Over?</title>
			<author>
				<first_name>Anno</first_name>
				<last_name>Birkin</last_name>
			</author>
			<genre>poetry</genre>
			<pages>80</pages>
			<ratings>
				<rating>5</rating>
				<rating>4</rating>
				<rating>4</rating>
				<rating>5</rating>
				<rating>3</rating>
			</ratings>
		</book>`
	count, err := library.AddBookXML([]byte(str))

	if err != nil {
		t.Fatalf("An error occured while parsing xml %s", err.Error())
	}

	if count != 1 {
		t.Errorf("Count must be 1 but is %d", count)
	}

	switch l := library.(type) {
	case *FancyLibrary:
		b := l.books["0954540018"]
		if b == nil {
			t.Fatalf("Book is nil, but it is expected to be in library")
		}

		if b.registeredCount != 1 {
			t.Errorf("Registered count of book should be 1, but is %s", b.registeredCount)
		}

		if b.availableCount != 1 {
			t.Errorf("Registered count of book should be 1, but is %s", b.availableCount)
		}

		if b.Author.FirstName != "Anno" {
			t.Errorf("Author's first name must be Anno, but is %s", b.Author.FirstName)
		}

		if b.Author.LastName != "Birkin" {
			t.Errorf("Author's last name must be Birkin, but is %s", b.Author.LastName)
		}

		if b.Title != "Who said the race is Over?" {
			t.Errorf("Title must be Who said the race is Over?, but is %s", b.Title)
		}
	}
}

func TestAddBook(t *testing.T) {
	library := NewLibrary(5)
	author := Person{FirstName: "Anno", LastName: "Birkin"}
	book := Book{ISBN: "0954540018", Author: author,
		registeredCount: 2, availableCount: 1,
		Title: "Who Said the Race is Over?"}
	switch l := library.(type) {
	case *FancyLibrary:
		l.addBook(&book)

		b := l.books["0954540018"]
		if b == nil {
			t.Errorf("Book is nil, but it is expected to be in library")
		}

		if b.registeredCount != 1 {
			t.Errorf("Registered count of book should be 1, but is %s", b.registeredCount)
		}

		if b.availableCount != 1 {
			t.Errorf("Registered count of book should be 1, but is %s", b.availableCount)
		}

		if b.Author.FirstName != "Anno" {
			t.Errorf("Author's first name must be Anno, but is %s", b.Author.FirstName)
		}

		if b.Author.LastName != "Birkin" {
			t.Errorf("Author's last name must be Birkin, but is %s", b.Author.LastName)
		}

		if b.Title != "Who Said the Race is Over?" {
			t.Errorf("Title must be Who Said the Race is Over?, but is %s", b.Title)
		}
	}
}

func TestTooManyBooks(t *testing.T) {
	library := NewLibrary(5)
	author := Person{FirstName: "Anno", LastName: "Birkin"}
	book := Book{ISBN: "0954540018", Author: author,
		registeredCount: 2, availableCount: 1,
		Title: "Who Said the Race is Over?"}

	author2 := Person{FirstName: "Anno", LastName: "Birkin"}
	book2 := Book{ISBN: "0954540018", Author: author2,
		registeredCount: 2, availableCount: 1,
		Title: "Who Said the Race is Over?"}

	author3 := Person{FirstName: "Anno", LastName: "Birkin"}
	book3 := Book{ISBN: "0954540018", Author: author3,
		registeredCount: 2, availableCount: 1,
		Title: "Who Said the Race is Over?"}

	author4 := Person{FirstName: "Anno", LastName: "Birkin"}
	book4 := Book{ISBN: "0954540018", Author: author4,
		registeredCount: 2, availableCount: 1,
		Title: "Who Said the Race is Over?"}

	author5 := Person{FirstName: "Anno", LastName: "Birkin"}
	book5 := Book{ISBN: "0954540018", Author: author5,
		registeredCount: 2, availableCount: 1,
		Title: "Who Said the Race is Over?"}

	switch l := library.(type) {
	case *FancyLibrary:
		l.addBook(&book)
		l.addBook(&book2)
		l.addBook(&book3)
		available, err := l.addBook(&book4)
		if available != 4 {
			t.Errorf("Available  must be 4, but is %d", available)
		}
		if err != nil {
			t.Errorf("Error must be nil, but is %s", err.Error())
		}
		available1, err1 := l.addBook(&book5)
		if available1 != 4 {
			t.Errorf("Available1  must be 4, but is %d", available1)
		}
		expectedError := "Има 4 копия на книга 0954540018"
		if err1 == nil {
			t.Errorf("Err1 is nil")
		} else if err1.Error() != expectedError {
			t.Errorf("Err1 is expected to be %s, but is %s", expectedError, err1.Error())
		}
	}
}

func TestHello(t *testing.T) {
	library := NewLibrary(5)
	request, response := library.Hello()
	if request == nil {
		t.Fatalf("Request must not be nil")
	}
	if response == nil {
		t.Fatalf("Response must not be nil")
	}
	switch l := library.(type) {
	case *FancyLibrary:
		select {
		case librarian := <-l.librarians:
			close(librarian.request)
		default:
			t.Fatalf("Channel with librarians has no content")
		}
	}
}

func TestLibrarianServe(t *testing.T) {
	library := NewLibrary(5)
	str := `
		<book isbn="0954540018">
			<title>Who said the race is Over?</title>
			<author>
				<first_name>Anno</first_name>
				<last_name>Birkin</last_name>
			</author>
			<genre>poetry</genre>
			<pages>80</pages>
			<ratings>
				<rating>5</rating>
				<rating>4</rating>
				<rating>4</rating>
				<rating>5</rating>
				<rating>3</rating>
			</ratings>
		</book>`
	library.AddBookXML([]byte(str))
	request, response := library.Hello()

	request <- &CoolLibraryRequest{"0954540018", BorrowBook}
	message := <-response
	book, err := message.GetBook()
	if err != nil {
		t.Errorf("There must not be an error")
	}

	if book == nil {
		t.Fatalf("Book must not be nil")
	}

	expected := "[0954540018] Who said the race is Over? от Anno Birkin"
	if book.String() != expected {
		t.Errorf("Expected book\n---\n%s\n---\nbut found\n---\n%s\n---\n", expected,
			book.String())
	}

	available, registered := message.GetAvailability()
	if available != 0 {
		t.Errorf("Expected available is 0, found %d", available)
	}
	if registered != 1 {
		t.Errorf("Expected registered is 1, found %d", registered)
	}

	request <- &CoolLibraryRequest{"0954540018", ReturnBook}
	message1 := <-response

	available1, registered1 := message1.GetAvailability()
	if available1 != 1 {
		t.Errorf("Expected available is 1, found %d", available1)
	}
	if registered1 != 1 {
		t.Errorf("Expected registered is 1, found %d", registered1)
	}

	close(request)
}

func TestLibrarianBorrowBook(t *testing.T) {
	request := make(chan LibraryRequest)
	response := make(chan LibraryResponse)
	library := NewLibrary(5)
	author := Person{FirstName: "Anno", LastName: "Birkin"}
	book := Book{ISBN: "0954540018", Author: author,
		registeredCount: 2, availableCount: 1,
		Title: "Who Said the Race is Over?"}

	switch l := library.(type) {
	case *FancyLibrary:
		l.books[book.ISBN] = &book
		librarian := Librarian{request: request, response: response, library: l}

		coolResponse := librarian.borrowBook("0954540018")
		if coolResponse.book.ISBN != "0954540018" {
			t.Errorf("The isbn should be 0954540018 but is %s", coolResponse.book.ISBN)
		}
		foundAvailable, foundRegistered := coolResponse.GetAvailability()
		if foundAvailable != 0 {
			t.Errorf("Available is expected to be 0, but is %d", foundAvailable)
		}
		if foundRegistered != 2 {
			t.Errorf("Registered is expected to be 2, but is %d", foundRegistered)
		}

		if coolResponse.err != nil {
			t.Errorf("response contains an error when is should not")
		}

		coolResponse1 := librarian.borrowBook("0954540018")

		if coolResponse1.err == nil {
			t.Fatalf("coolResponse1 must have an error but there isn't one")
		}

		expectedError := "Няма наличност на книга 0954540018"
		if coolResponse1.err.Error() != expectedError {
			if coolResponse1.err.Error() != expectedError {
				t.Errorf("Expected error\n---\n%s\n---\nbut found\n---\n%s\n---\n", expectedError,
					coolResponse1.err.Error())
			}
		}

		coolResponse2 := librarian.borrowBook("IMNOTABOOK")
		if coolResponse2.err == nil {
			t.Fatalf("coolResponse1 must have an error but there isn't one")
		}

		expectedError2 := "Непозната книга IMNOTABOOK"
		if coolResponse2.err.Error() != expectedError2 {
			t.Errorf("Expected error\n---\n%s\n---\nbut found\n---\n%s\n---\n", expectedError2,
				coolResponse2.err.Error())
		}
	}
}

func TestLibrarianReturnBook(t *testing.T) {
	request := make(chan LibraryRequest)
	response := make(chan LibraryResponse)
	library := NewLibrary(5)
	author := Person{FirstName: "Anno", LastName: "Birkin"}
	book := Book{ISBN: "0954540018", Author: author,
		registeredCount: 2, availableCount: 1,
		Title: "Who Said the Race is Over?"}

	switch l := library.(type) {
	case *FancyLibrary:
		l.books[book.ISBN] = &book
		librarian := Librarian{request: request, response: response, library: l}
		coolResponse := librarian.returnBook("0954540018")
		foundAvailable, foundRegistered := coolResponse.GetAvailability()

		if foundAvailable != 2 {
			t.Errorf("Found Available must be 2, but is %d", foundAvailable)
		}

		if foundRegistered != 2 {
			t.Errorf("Found Registered must be 2, but is %d", foundRegistered)
		}

		coolResponse1 := librarian.returnBook("0954540018")
		if coolResponse1.err == nil {
			t.Fatalf("There must be an error, but there is not")
		}

		expectedError := "Всички копия са налични 0954540018"
		if coolResponse1.err.Error() != expectedError {
			t.Errorf("Expected error\n---\n%s\n---\nbut found\n---\n%s\n---\n", expectedError,
				coolResponse1.err.Error())
		}

		coolResponse2 := librarian.returnBook("IMNOTABOOK")
		if coolResponse2.err == nil {
			t.Fatalf("There must be an error, but there is not")
		}

		expectedError2 := "Непозната книга IMNOTABOOK"
		if coolResponse2.err.Error() != expectedError2 {
			t.Errorf("Expected error\n---\n%s\n---\nbut found\n---\n%s\n---\n", expectedError2,
				coolResponse2.err.Error())
		}
	}
}

func TestLibrarianGetAvailability(t *testing.T) {
	request := make(chan LibraryRequest)
	response := make(chan LibraryResponse)
	library := NewLibrary(5)
	author := Person{FirstName: "Anno", LastName: "Birkin"}
	book := Book{ISBN: "0954540018", Author: author,
		registeredCount: 2, availableCount: 1,
		Title: "Who Said the Race is Over?"}

	switch l := library.(type) {
	case *FancyLibrary:
		l.books[book.ISBN] = &book
		librarian := Librarian{request: request, response: response, library: l}
		coolResponse := librarian.getAvailability("0954540018")

		foundAvailable, foundRegistered := coolResponse.GetAvailability()

		if foundAvailable != 1 {
			t.Errorf("Found Available must be 1, but is %d", foundAvailable)
		}

		if foundRegistered != 2 {
			t.Errorf("Found Registered must be 2, but is %d", foundRegistered)
		}
		expected := "[0954540018] Who Said the Race is Over? от Anno Birkin"
		found, err := coolResponse.GetBook()
		if found.String() != expected {
			t.Errorf("Expected\n---\n%s\n---\nbut found\n---\n%s\n---\n", expected, found.String())
		}

		if err != nil {
			t.Errorf("There must not be an error")
		}

		coolResponse1 := librarian.getAvailability("IMNOTABOOK")

		found1, err1 := coolResponse1.GetBook()
		if found1 != nil {
			t.Errorf("Content must be nil")
		}

		if err1 == nil {
			t.Fatalf("There must be an error")
		}

		if err1.Error() != "Непозната книга IMNOTABOOK" {
			t.Errorf("Expected error %s, found %s", "Непозната книга IMNOTABOOK", err1.Error())
		}
	}
}

func TestMoreThanOneReqest(t *testing.T) {
	library := NewLibrary(5)
	str := `
		<book isbn="0954540018">
			<title>Who said the race is Over?</title>
			<author>
				<first_name>Anno</first_name>
				<last_name>Birkin</last_name>
			</author>
			<genre>poetry</genre>
			<pages>80</pages>
			<ratings>
				<rating>5</rating>
				<rating>4</rating>
				<rating>4</rating>
				<rating>5</rating>
				<rating>3</rating>
			</ratings>
		</book>`
	library.AddBookXML([]byte(str))
	request, response := library.Hello()

	request <- &CoolLibraryRequest{"0954540018", BorrowBook}
	request <- &CoolLibraryRequest{"0954540018", ReturnBook}

	message := <-response
	book, err := message.GetBook()
	if err != nil {
		t.Errorf("There must not be an error")
	}

	if book == nil {
		t.Fatalf("Book must not be nil")
	}

	expected := "[0954540018] Who said the race is Over? от Anno Birkin"
	if book.String() != expected {
		t.Errorf("Expected book\n---\n%s\n---\nbut found\n---\n%s\n---\n", expected,
			book.String())
	}

	available, registered := message.GetAvailability()
	if available != 1 {
		t.Errorf("Expected available is 1, found %d", available)
	}
	if registered != 1 {
		t.Errorf("Expected registered is 1, found %d", registered)
	}

	message1 := <-response

	available1, registered1 := message1.GetAvailability()
	if available1 != 1 {
		t.Errorf("Expected available is 1, found %d", available1)
	}
	if registered1 != 1 {
		t.Errorf("Expected registered is 1, found %d", registered1)
	}

	close(request)
}

func TestMoreThanOneLibrarian(t *testing.T) {
	library := NewLibrary(5)
	str := `
		<book isbn="0954540018">
			<title>Who said the race is Over?</title>
			<author>
				<first_name>Anno</first_name>
				<last_name>Birkin</last_name>
			</author>
			<genre>poetry</genre>
			<pages>80</pages>
			<ratings>
				<rating>5</rating>
				<rating>4</rating>
				<rating>4</rating>
				<rating>5</rating>
				<rating>3</rating>
			</ratings>
		</book>`
	library.AddBookXML([]byte(str))
	request1, response1 := library.Hello()
	request2, response2 := library.Hello()

	request1 <- &CoolLibraryRequest{"0954540018", BorrowBook}
	message1 := <-response1
	_, err := message1.GetBook()
	available1, _ := message1.GetAvailability()

	if err != nil {
		t.Errorf("There must not be an error")
	}

	if available1 != 0 {
		t.Errorf("Available is expected to be 0, but is %d", available1)
	}

	request2 <- &CoolLibraryRequest{"0954540018", BorrowBook}
	message2 := <-response2
	_, err2 := message2.GetBook()

	if err2 == nil {
		t.Fatalf("There must be an error")
	}

	close(request1)
	close(request2)
}

func TestAddTwoBooks(t *testing.T) {
	library := NewLibrary(5)
	str := `
		<book isbn="0954540018">
			<title>Who said the race is Over?</title>
			<author>
				<first_name>Anno</first_name>
				<last_name>Birkin</last_name>
			</author>
			<genre>poetry</genre>
			<pages>80</pages>
			<ratings>
				<rating>5</rating>
				<rating>4</rating>
				<rating>4</rating>
				<rating>5</rating>
				<rating>3</rating>
			</ratings>
		</book>`
	library.AddBookXML([]byte(str))
	library.AddBookXML([]byte(str))

	request, response := library.Hello()
	request <- &CoolLibraryRequest{"0954540018", GetAvailability}

	message := <-response

	available, registered := message.GetAvailability()
	if available != 2 {
		t.Errorf("available books must be 2, but are %d", available)
	}

	if registered != 2 {
		t.Errorf("registered books must be 2, but are %d", registered)
	}
}

// no error messages in this test
// just make sure there is no deadlock
func TestManyLibrarians(t *testing.T) {
	library := NewLibrary(5)
	str := `
		<book isbn="0954540018">
			<title>Who said the race is Over?</title>
			<author>
				<first_name>Anno</first_name>
				<last_name>Birkin</last_name>
			</author>
			<genre>poetry</genre>
			<pages>80</pages>
			<ratings>
				<rating>5</rating>
				<rating>4</rating>
				<rating>4</rating>
				<rating>5</rating>
				<rating>3</rating>
			</ratings>
		</book>`
	library.AddBookXML([]byte(str))
	library.AddBookXML([]byte(str))

	go func() {
		for i := 0; i < 100; i++ {
			request, response := library.Hello()
			for j := 0; j < 1000000; j++ {
				request <- &CoolLibraryRequest{"0954540018", GetAvailability}
				<-response
			}
			close(request)
		}
	}()
}
