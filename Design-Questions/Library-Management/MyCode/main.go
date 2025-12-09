package main

import "time"
type BookStatus string

var (
	Available BookStatus = "available"
	Taken      BookStatus = "taken"
)

func (b BookStatus) ToString() string {
	return string(b)
}

type TransactionStatus string

var (
	Pending TransactionStatus = "pending"
	Success      TransactionStatus = "success"
)

func (b TransactionStatus) ToString() string {
	return string(b)
}

type User struct {
	UserId, name, email, mobile string
	books                       []*Book
	date                        []*Date
	Fines                       []*Transaction
}

func NewUser(u, n, e, m string) User {
	return User{
		UserId: u, name: n, email: e, mobile: m,
		books: make([]*Book, 0),
		date:  make([]*Date, 0),
		Fines: make([]*Transaction, 0),
	}
}

func (u *User) TakeBook(book Book) bool {
	if book.Status == Taken || !u.IsUserEligible() {
		return false
	} else {
		bookDate := NewDate(time.Now(), book.BookId)
		u.date = append(u.date, &bookDate)
		u.books = append(u.books, &book)
		book.ChangeBookStatus(Taken)
		return true

	}

}

func (u *User) IsUserEligible() bool {
	for _, val := range u.Fines {
		if val.Status == Pending {
			return false
		}
	}
	return true
}

// greater than 7 days fine each day 10 rupees
func (u *User) ReturnBook(book Book) bool {
	for _, d := range u.date {
		if d.BookId == book.BookId {
			takenDate := d.AllotedDate
			diff := time.Now().Sub(takenDate)
			if diff > 7 {
				txn := NewTransaction("txn123", book.BookId, int64(diff-7)*10)
				u.Fines = append(u.Fines, &txn)
			}
			return true
		}
	}
	return false
}

type Library struct{
	Users []User
	Books []Book

}

func NewLibrary()Library{
	return Library{
		Users: make([]User, 0),
		Books: make([]Book, 0),
	}
}

func(l *Library)AddBook(book ...Book){
	for _,b:=range book{
		l.Books=append(l.Books, b)
	}
}

func(l *Library)AddUser(user ...User){
	for _,b:=range user{
		l.Users=append(l.Users, b)
	}
}

func(l *Library)AllotBookToUser(user User,book Book)bool{
	return user.TakeBook(book)
}

func(l *Library)ReturnBook(user User,book Book)bool{
	return user.ReturnBook(book)
}

type Date struct {
	AllotedDate time.Time
	ReturnDate  time.Time
	BookId      string
}

type Transaction struct {
	TransactionId string
	BookId        string
	Fine          int64
	Status        TransactionStatus
}

func NewDate(allotedDate time.Time, bookId string) Date {
	return Date{
		AllotedDate: allotedDate,
		ReturnDate:  time.Time{},
		BookId:      bookId,
	}
}

func NewTransaction(txnId, bId string, fine int64) Transaction {
	return Transaction{
		TransactionId: txnId,
		BookId:        bId,
		Fine:          fine,
		Status:        Pending,
	}
}

func (t *Transaction) PayFine(txnId, bookId string) {
	if t.BookId == bookId && t.TransactionId == txnId && t.Status.ToString() == "pending" {
		t.Status = Success
	}
}

type Book struct {
	BookId       string
	Name, Author string
	Status       BookStatus
}

func NewBook(id, name, author string) Book {
	return Book{BookId: id, Name: name, Author: author, Status: Available}
}

func (b *Book) BookStatus() string {
	return b.Status.ToString()
}

func(b *Book)ChangeBookStatus(status BookStatus){
	b.Status=status
}