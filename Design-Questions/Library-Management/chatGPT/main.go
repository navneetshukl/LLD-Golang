package main

import (
	"errors"
	"fmt"
	"time"
)

//
// ---------- ENUMS ----------
//

type BookStatus string

const (
	BookAvailable BookStatus = "available"
	BookTaken     BookStatus = "taken"
)

type TransactionStatus string

const (
	TxnPending TransactionStatus = "pending"
	TxnSuccess TransactionStatus = "success"
)

//
// ---------- DOMAIN ENTITIES ----------
//

type Book struct {
	ID     string
	Title  string
	Author string
	Status BookStatus
}

func NewBook(id, title, author string) *Book {
	return &Book{
		ID:     id,
		Title:  title,
		Author: author,
		Status: BookAvailable,
	}
}

func (b *Book) IsAvailable() bool {
	return b.Status == BookAvailable
}

func (b *Book) ChangeStatus(status BookStatus) {
	b.Status = status
}

type BorrowRecord struct {
	BookID     string
	BorrowedAt time.Time
	ReturnedAt *time.Time
}

//
// ---------- TRANSACTION (Now Stored in Library) ----------
//

type Transaction struct {
	ID        string
	UserID    string
	BookID    string
	Fine      int64
	Status    TransactionStatus
	CreatedAt time.Time
	PaidAt    *time.Time
}

func NewTransaction(id, userID, bookID string, fine int64) *Transaction {
	return &Transaction{
		ID:        id,
		UserID:    userID,
		BookID:    bookID,
		Fine:      fine,
		Status:    TxnPending,
		CreatedAt: time.Now(),
	}
}

func (t *Transaction) MarkPaid() {
	if t.Status == TxnPending {
		t.Status = TxnSuccess
		now := time.Now()
		t.PaidAt = &now
	}
}

//
// ---------- USER (Now Clean — No Transactions Stored Inside) ----------
//

type User struct {
	ID       string
	Name     string
	Email    string
	Mobile   string
	Borrowed map[string]*BorrowRecord // key: bookID
}

func NewUser(id, name, email, mobile string) *User {
	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		Mobile:   mobile,
		Borrowed: make(map[string]*BorrowRecord),
	}
}

// This checks if Library has pending fines for this user.
func (u *User) GetActiveBorrow(bookID string) *BorrowRecord {
	br, ok := u.Borrowed[bookID]
	if !ok {
		return nil
	}
	if br.ReturnedAt != nil {
		return nil
	}
	return br
}

//
// ---------- LIBRARY (AGGREGATE ROOT) ----------
//

type Library struct {
	users        map[string]*User
	books        map[string]*Book
	transactions map[string]*Transaction // NEW: All transactions stored centrally
}

func NewLibrary() *Library {
	return &Library{
		users:        make(map[string]*User),
		books:        make(map[string]*Book),
		transactions: make(map[string]*Transaction),
	}
}

func (l *Library) AddUser(u *User) { l.users[u.ID] = u }
func (l *Library) AddBook(b *Book) { l.books[b.ID] = b }

func (l *Library) GetUser(id string) (*User, error) {
	u, ok := l.users[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (l *Library) GetBook(id string) (*Book, error) {
	b, ok := l.books[id]
	if !ok {
		return nil, errors.New("book not found")
	}
	return b, nil
}

//
// ---------- HELPER: Check pending fine inside Library ----------
//

func (l *Library) UserHasPendingFine(userID string) bool {
	for _, txn := range l.transactions {
		if txn.UserID == userID && txn.Status == TxnPending {
			return true
		}
	}
	return false
}

//
// ---------- BORROW BOOK ----------
//

func (l *Library) BorrowBook(userID, bookID string) error {
	user, err := l.GetUser(userID)
	if err != nil {
		return err
	}

	book, err := l.GetBook(bookID)
	if err != nil {
		return err
	}

	if !book.IsAvailable() {
		return errors.New("book not available")
	}

	if l.UserHasPendingFine(userID) {
		return errors.New("user has pending fines")
	}

	user.Borrowed[bookID] = &BorrowRecord{
		BookID:     bookID,
		BorrowedAt: time.Now(),
	}

	book.ChangeStatus(BookTaken)
	return nil
}

//
// ---------- RETURN BOOK ----------
//

func (l *Library) ReturnBook(userID, bookID string) (*Transaction, error) {
	user, err := l.GetUser(userID)
	if err != nil {
		return nil, err
	}

	book, err := l.GetBook(bookID)
	if err != nil {
		return nil, err
	}

	active := user.GetActiveBorrow(bookID)
	if active == nil {
		return nil, errors.New("no active borrow found")
	}

	now := time.Now()
	active.ReturnedAt = &now
	book.ChangeStatus(BookAvailable)

	// Fine rule: after 7 days → 10 per extra day
	days := int(now.Sub(active.BorrowedAt).Hours() / 24)

	if days <= 7 {
		return nil, nil // no fine
	}

	extraDays := days - 7
	fineAmount := int64(extraDays * 10)

	txnID := fmt.Sprintf("TXN-%d", time.Now().UnixNano())
	txn := NewTransaction(txnID, userID, bookID, fineAmount)

	l.transactions[txnID] = txn

	return txn, nil
}

//
// ---------- PAY FINE ----------
//

func (l *Library) PayFine(userID, txnID string) error {
	txn, ok := l.transactions[txnID]
	if !ok {
		return errors.New("transaction not found")
	}

	if txn.UserID != userID {
		return errors.New("transaction does not belong to this user")
	}

	if txn.Status == TxnSuccess {
		return errors.New("fine already paid")
	}

	txn.MarkPaid()
	return nil
}

//
// ---------- DEMO MAIN ----------
//

func main() {
	library := NewLibrary()

	user := NewUser("U1", "Navneet", "navneet@example.com", "9999999999")
	book1 := NewBook("B1", "Clean Code", "Robert C. Martin")

	library.AddUser(user)
	library.AddBook(book1)

	// Borrow
	if err := library.BorrowBook("U1", "B1"); err != nil {
		fmt.Println("Borrow error:", err)
		return
	}
	fmt.Println("Book borrowed")

	// Return (fine may or may not be generated)
	txn, err := library.ReturnBook("U1", "B1")
	if err != nil {
		fmt.Println("Return error:", err)
		return
	}

	if txn == nil {
		fmt.Println("Returned with no fine")
		return
	}

	fmt.Printf("Fine generated: %d, TxnID=%s\n", txn.Fine, txn.ID)

	// Pay fine
	if err := library.PayFine("U1", txn.ID); err != nil {
		fmt.Println("Pay error:", err)
		return
	}

	fmt.Println("Fine paid")
}
