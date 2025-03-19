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

	provider preloader.LoadableProvider
}

func (u *User) GetResourceID() UserID {
	return u.ID
}

func (u *User) Books(ctx context.Context) ([]*Book, error) {
	loadable := preloader.GetLoadableFromProvider[UserID, *User, BookID, *Book](u.provider, "Books")
	return loadable.Load(ctx, u)
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
	PlaceID  PlaceID
	
	provider preloader.LoadableProvider
}

func (u *Book) GetResourceID() BookID {
	return u.ID
}

func (b *Book) Author(ctx context.Context) (*User, error) {
	loadable := preloader.GetHasOneLoadableFromProvider[BookID, *Book, UserID, *User](b.provider, "Authors")
	return loadable.Load(ctx, b)
}

func (b *Book) Place(ctx context.Context) (*Place, error) {
	loadable := preloader.GetHasOneLoadableFromProvider[BookID, *Book, PlaceID, *Place](b.provider, "Places")
	return loadable.Load(ctx, b)
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
	
	provider preloader.LoadableProvider
}

func (p *Place) GetResourceID() PlaceID {
	return p.ID
}
