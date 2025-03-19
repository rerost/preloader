package preloader

import "context"

// TypedResource is a type-safe version of Resource that enforces loadable registration at compile time
type TypedResource[TID comparable] interface {
	Resource[TID]
	SetProvider(provider *TypedLoadableProvider)
}

// Example type definitions for demonstration purposes
type UserID int
type BookID int
type PlaceID int

// TypedUser is an example of a type-safe user model
type TypedUser struct {
	ID   UserID
	Name string
	
	provider *TypedLoadableProvider
	// Registered loadables
	booksLoadable RegisteredLoadable[Registered, UserID, *TypedUser, BookID, *TypedBook]
}

func (u *TypedUser) GetResourceID() UserID {
	return u.ID
}

func (u *TypedUser) SetProvider(provider *TypedLoadableProvider) {
	u.provider = provider
}

// RegisterBooksLoadable registers the Books loadable for this user
func (u *TypedUser) RegisterBooksLoadable(
	loadable RegisteredLoadable[Registered, UserID, *TypedUser, BookID, *TypedBook],
) {
	u.booksLoadable = loadable
}

// Books returns the books for this user
// This will fail at compile time if the "Books" loadable is not registered
func (u *TypedUser) Books(ctx context.Context) ([]*TypedBook, error) {
	loadable := GetRegisteredLoadable(u.provider, u.booksLoadable)
	return loadable.Load(ctx, u)
}

// TypedBook is an example of a type-safe book model
type TypedBook struct {
	ID    BookID
	Title string
	
	AuthorID UserID
	PlaceID  PlaceID
	
	provider *TypedLoadableProvider
	// Registered loadables
	authorLoadable RegisteredHasOneLoadable[Registered, BookID, *TypedBook, UserID, *TypedUser]
	placeLoadable  RegisteredHasOneLoadable[Registered, BookID, *TypedBook, PlaceID, *TypedPlace]
}

func (b *TypedBook) GetResourceID() BookID {
	return b.ID
}

func (b *TypedBook) SetProvider(provider *TypedLoadableProvider) {
	b.provider = provider
}

// RegisterAuthorLoadable registers the Author loadable for this book
func (b *TypedBook) RegisterAuthorLoadable(
	loadable RegisteredHasOneLoadable[Registered, BookID, *TypedBook, UserID, *TypedUser],
) {
	b.authorLoadable = loadable
}

// RegisterPlaceLoadable registers the Place loadable for this book
func (b *TypedBook) RegisterPlaceLoadable(
	loadable RegisteredHasOneLoadable[Registered, BookID, *TypedBook, PlaceID, *TypedPlace],
) {
	b.placeLoadable = loadable
}

// Author returns the author of this book
// This will fail at compile time if the "Authors" loadable is not registered
func (b *TypedBook) Author(ctx context.Context) (*TypedUser, error) {
	loadable := GetRegisteredHasOneLoadable(b.provider, b.authorLoadable)
	return loadable.Load(ctx, b)
}

// Place returns the place of this book
// This will fail at compile time if the "Places" loadable is not registered
func (b *TypedBook) Place(ctx context.Context) (*TypedPlace, error) {
	loadable := GetRegisteredHasOneLoadable(b.provider, b.placeLoadable)
	return loadable.Load(ctx, b)
}

// TypedPlace is an example of a type-safe place model
type TypedPlace struct {
	ID   PlaceID
	Name string
	
	provider *TypedLoadableProvider
}

func (p *TypedPlace) GetResourceID() PlaceID {
	return p.ID
}

func (p *TypedPlace) SetProvider(provider *TypedLoadableProvider) {
	p.provider = provider
}
