package main

import (
	"context"
	"fmt"

	"github.com/rerost/preloader"
)

// TypedUser is a type-safe version of User
type TypedUser struct {
	ID   UserID
	Name string
	
	provider *preloader.TypedLoadableProvider
	booksLoadable preloader.RegisteredLoadable[preloader.Registered, UserID, *TypedUser, BookID, *TypedBook]
}

func (u *TypedUser) GetResourceID() UserID {
	return u.ID
}

func (u *TypedUser) SetProvider(provider *preloader.TypedLoadableProvider) {
	u.provider = provider
}

// RegisterBooksLoadable registers the Books loadable for this user
func (u *TypedUser) RegisterBooksLoadable(
	loadable preloader.RegisteredLoadable[preloader.Registered, UserID, *TypedUser, BookID, *TypedBook],
) {
	u.booksLoadable = loadable
}

// Books returns the books for this user
func (u *TypedUser) Books(ctx context.Context) ([]*TypedBook, error) {
	loadable := preloader.GetRegisteredLoadable(u.provider, u.booksLoadable)
	return loadable.Load(ctx, u)
}

// TypedBook is a type-safe version of Book
type TypedBook struct {
	ID    BookID
	Title string
	
	AuthorID UserID
	PlaceID  PlaceID
	
	provider *preloader.TypedLoadableProvider
	authorLoadable preloader.RegisteredHasOneLoadable[preloader.Registered, BookID, *TypedBook, UserID, *TypedUser]
	placeLoadable  preloader.RegisteredHasOneLoadable[preloader.Registered, BookID, *TypedBook, PlaceID, *Place]
}

func (b *TypedBook) GetResourceID() BookID {
	return b.ID
}

func (b *TypedBook) SetProvider(provider *preloader.TypedLoadableProvider) {
	b.provider = provider
}

// RegisterAuthorLoadable registers the Author loadable for this book
func (b *TypedBook) RegisterAuthorLoadable(
	loadable preloader.RegisteredHasOneLoadable[preloader.Registered, BookID, *TypedBook, UserID, *TypedUser],
) {
	b.authorLoadable = loadable
}

// RegisterPlaceLoadable registers the Place loadable for this book
func (b *TypedBook) RegisterPlaceLoadable(
	loadable preloader.RegisteredHasOneLoadable[preloader.Registered, BookID, *TypedBook, PlaceID, *Place],
) {
	b.placeLoadable = loadable
}

// Author returns the author of this book
func (b *TypedBook) Author(ctx context.Context) (*TypedUser, error) {
	loadable := preloader.GetRegisteredHasOneLoadable(b.provider, b.authorLoadable)
	return loadable.Load(ctx, b)
}

// Place returns the place of this book
func (b *TypedBook) Place(ctx context.Context) (*Place, error) {
	loadable := preloader.GetRegisteredHasOneLoadable(b.provider, b.placeLoadable)
	return loadable.Load(ctx, b)
}

// TypedUserRepository is a type-safe version of UserRepository
type TypedUserRepository struct {
	m        map[UserID]*TypedUser
	provider *preloader.TypedLoadableProvider
}

func NewTypedUserRepository(provider *preloader.TypedLoadableProvider) *TypedUserRepository {
	repo := &TypedUserRepository{
		m:        make(map[UserID]*TypedUser),
		provider: provider,
	}
	
	// Initialize users
	repo.m[1] = &TypedUser{
		ID:       1,
		Name:     "Alice",
		provider: provider,
	}
	repo.m[2] = &TypedUser{
		ID:       2,
		Name:     "Bob",
		provider: provider,
	}
	
	return repo
}

func (u *TypedUserRepository) List(ctx context.Context, ids []UserID) ([]*TypedUser, error) {
	users := make([]*TypedUser, len(ids))
	for i, id := range ids {
		users[i] = &TypedUser{
			ID:       id,
			Name:     fmt.Sprintf("User %d", id),
			provider: u.provider,
		}
	}
	return users, nil
}

func (u *TypedUserRepository) All() ([]*TypedUser, error) {
	users := make([]*TypedUser, 0, len(u.m))
	for _, user := range u.m {
		users = append(users, user)
	}
	return users, nil
}

// Example of using the type-safe API
func typeSafeExample() {
	// Create a type-safe provider
	provider := preloader.NewTypedLoadableProvider()
	
	// Create and register loadables
	bookLoadable := preloader.EmptyLoadable[UserID, *TypedUser, BookID, *TypedBook]()
	registeredBookLoadable := preloader.RegisterLoadable(provider, preloader.LoadableKey("Books"), bookLoadable)
	
	authorLoadable := preloader.EmptyHasOneLoadable[BookID, *TypedBook, UserID, *TypedUser]()
	registeredAuthorLoadable := preloader.RegisterHasOneLoadable(provider, preloader.LoadableKey("Authors"), authorLoadable)
	
	// Create a user and register the loadable
	user := &TypedUser{ID: 1, Name: "Test User"}
	user.SetProvider(provider)
	user.RegisterBooksLoadable(registeredBookLoadable)
	
	// This will compile because the loadable is registered
	books, _ := user.Books(context.Background())
	fmt.Println("Books:", books)
	
	// Create a book and register the loadable
	book := &TypedBook{ID: 1, Title: "Test Book"}
	book.SetProvider(provider)
	book.RegisterAuthorLoadable(registeredAuthorLoadable)
	
	// This will compile because the loadable is registered
	author, _ := book.Author(context.Background())
	fmt.Println("Author:", author)
}
