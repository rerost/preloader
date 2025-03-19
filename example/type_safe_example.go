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
}

func (u *TypedUser) GetResourceID() UserID {
	return u.ID
}

func (u *TypedUser) SetProvider(provider *preloader.TypedLoadableProvider) {
	u.provider = provider
}

// Books returns the books for this user
func (u *TypedUser) Books(ctx context.Context) ([]*TypedBook, error) {
	loadable := preloader.MustGetLoadable[UserID, *TypedUser, BookID, *TypedBook](u.provider, preloader.LoadableKey("Books"))
	return loadable.Load(ctx, u)
}

// TypedBook is a type-safe version of Book
type TypedBook struct {
	ID    BookID
	Title string
	
	AuthorID UserID
	PlaceID  PlaceID
	
	provider *preloader.TypedLoadableProvider
}

func (b *TypedBook) GetResourceID() BookID {
	return b.ID
}

func (b *TypedBook) SetProvider(provider *preloader.TypedLoadableProvider) {
	b.provider = provider
}

// Author returns the author of this book
func (b *TypedBook) Author(ctx context.Context) (*TypedUser, error) {
	loadable := preloader.MustGetHasOneLoadable[BookID, *TypedBook, UserID, *TypedUser](b.provider, preloader.LoadableKey("Authors"))
	return loadable.Load(ctx, b)
}

// Place returns the place of this book
func (b *TypedBook) Place(ctx context.Context) (*Place, error) {
	loadable := preloader.MustGetHasOneLoadable[BookID, *TypedBook, PlaceID, *Place](b.provider, preloader.LoadableKey("Places"))
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
	// If we forget to register a loadable, we'll get a compile-time error
	bookLoadable := preloader.EmptyLoadable[UserID, *TypedUser, BookID, *TypedBook]()
	preloader.RegisterTypedLoadable(provider, preloader.LoadableKey("Books"), bookLoadable)
	
	// This would cause a compile-time error if we try to use Author() without registering the Authors loadable
	authorLoadable := preloader.EmptyHasOneLoadable[BookID, *TypedBook, UserID, *TypedUser]()
	preloader.RegisterTypedHasOneLoadable(provider, preloader.LoadableKey("Authors"), authorLoadable)
}
