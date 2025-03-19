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
}

func (u *TypedUser) GetResourceID() UserID {
	return u.ID
}

func (u *TypedUser) SetProvider(provider *TypedLoadableProvider) {
	u.provider = provider
}

// Books returns the books for this user
// This will fail at compile time if the "Books" loadable is not registered
func (u *TypedUser) Books(ctx context.Context) ([]*TypedBook, error) {
	loadable := MustGetLoadable[UserID, *TypedUser, BookID, *TypedBook](u.provider, LoadableKey("Books"))
	return loadable.Load(ctx, u)
}

// TypedBook is an example of a type-safe book model
type TypedBook struct {
	ID    BookID
	Title string
	
	AuthorID UserID
	PlaceID  PlaceID
	
	provider *TypedLoadableProvider
}

func (b *TypedBook) GetResourceID() BookID {
	return b.ID
}

func (b *TypedBook) SetProvider(provider *TypedLoadableProvider) {
	b.provider = provider
}

// Author returns the author of this book
// This will fail at compile time if the "Authors" loadable is not registered
func (b *TypedBook) Author(ctx context.Context) (*TypedUser, error) {
	loadable := MustGetHasOneLoadable[BookID, *TypedBook, UserID, *TypedUser](b.provider, LoadableKey("Authors"))
	return loadable.Load(ctx, b)
}

// Place returns the place of this book
// This will fail at compile time if the "Places" loadable is not registered
func (b *TypedBook) Place(ctx context.Context) (*TypedPlace, error) {
	loadable := MustGetHasOneLoadable[BookID, *TypedBook, PlaceID, *TypedPlace](b.provider, LoadableKey("Places"))
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
