package main

import "strconv"

type Book struct {
	ID          int
	Title       string
	Author      string
	IsAvailable bool
}

type Magazine struct {
	ID          int
	Title       string
	Issue       int
	IsAvailable bool
}

// Library结构体定义
type Library struct {
	Books     []*Book
	Magazines []*Magazine
	Name      string
}

type Manageable interface {
	Borrow() bool
	Return() bool
	GetInfo() string
}

func main() {
	// 创建图书馆实例
	library := Library{
		Name: "Library",
	}

	// 添加书籍和杂志到图书馆
	library.AddBook(&Book{ID: 1, Title: "Go Programming", Author: "John Doe", IsAvailable: true})
	library.AddBook(&Book{ID: 2, Title: "Advanced Go", Author: "Jane Smith", IsAvailable: false})
	library.AddBook(&Book{ID: 3, Title: "Advanced Go", Author: "Jane Smith", IsAvailable: false})
	library.AddMagazine(&Magazine{ID: 1, Title: "Tech Today", Issue: 42, IsAvailable: true})
	library.AddMagazine(&Magazine{ID: 2, Title: "Dev Monthly", Issue: 15, IsAvailable: true})
	library.AddMagazine(&Magazine{ID: 3, Title: "Dev Monthly", Issue: 16, IsAvailable: false})
	if book := library.FindBookByID(1); book != nil {
		book.Borrow()
	}
	library.PrintAvailableItems()
}

func (b *Book) GetInfo() string {
	return "Book ID: " + strconv.Itoa(b.ID) + ", Title: " + b.Title + ", Author: " + b.Author + ", Available: " + strconv.FormatBool(b.IsAvailable)
}
func (m *Magazine) GetInfo() string {
	return "Book ID: " + strconv.Itoa(m.ID) + ", Title: " + m.Title + ", Issue: " + strconv.Itoa(m.Issue) + ", Available: " + strconv.FormatBool(m.IsAvailable)
}

func (b *Book) Borrow() bool {
	if b.IsAvailable {
		b.IsAvailable = false
		return true
	}
	return false
}
func (b *Book) Return() bool {
	if !b.IsAvailable {
		b.IsAvailable = true
		return true
	}
	return false
}

func (l *Library) AddBook(book *Book) {
	l.Books = append(l.Books, book)
}

func (l *Library) FindBookByID(id int) *Book {
	for _, book := range l.Books {
		if book.ID == id {
			return book
		}
	}
	return nil
}

func (l *Library) GetAvailableBooks() []*Book {
	var available []*Book
	for _, book := range l.Books {
		if book.IsAvailable {
			available = append(available, book)
		}
	}
	return available
}

func (l *Library) AddMagazine(magazine *Magazine) {
	l.Magazines = append(l.Magazines, magazine)
}

func (l *Library) FindMagazineByID(id int) *Magazine {
	for _, mag := range l.Magazines {
		if mag.ID == id {
			return mag
		}
	}
	return nil
}

func (l *Library) GetAvailableMagazines() []*Magazine {
	var available []*Magazine
	for _, mag := range l.Magazines {
		if mag.IsAvailable {
			available = append(available, mag)
		}
	}
	return available
}

func (l *Library) PrintAvailableItems() {
	println("\nAvailable Books:")
	for _, book := range l.GetAvailableBooks() {
		println(book.GetInfo())
	}

	println("\nAvailable Magazines:")
	for _, mag := range l.GetAvailableMagazines() {
		println(mag.GetInfo())
	}
}
