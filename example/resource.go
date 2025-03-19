package main

import (
	"fmt"

	"github.com/rerost/preloader"
)

// User
type UserID int

func (u *UserID) String() string {
	return fmt.Sprintf("User %d", *u)
}

type User struct {
	ID   UserID
	Name string

	Books preloader.Loadable[UserID, *User, BookID, *Book]
}

func (u *User) GetResourceID() UserID {
	return u.ID
}

// Book
type BookID int

func (b *BookID) String() string {
	return fmt.Sprintf("Book %d", *b)
}

type Book struct {
	ID    BookID
	Title string

	AuthorID UserID
	Author   preloader.HasOneLoadable[BookID, *Book, UserID, *User]

	PlaceID PlaceID
	Place   preloader.HasOneLoadable[BookID, *Book, PlaceID, *Place]
}

func (u *Book) GetResourceID() BookID {
	return u.ID
}

// Place
type PlaceID int

func (p *PlaceID) String() string {
	return fmt.Sprintf("Place %d", *p)
}

type Place struct {
	ID   PlaceID
	Name string
	Type string
}

func (p *Place) GetResourceID() PlaceID {
	return p.ID
}
