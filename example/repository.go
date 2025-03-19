package main

import (
	"context"
	"fmt"

	"github.com/rerost/preloader"
)

type UserRepository interface {
	List(ctx context.Context, id []UserID) ([]*User, error)
	All() ([]*User, error)
}

type BookRepository interface {
	List(ctx context.Context, id []BookID) ([]*Book, error)
	ByUsers(ctx context.Context, userIDs []UserID) (map[UserID][]BookID, error)
}

type PlaceRepository interface {
	List(ctx context.Context, ids []PlaceID) ([]*Place, error)
}

func NewUserRepository(
	booksLoadable preloader.Loadable[UserID, *User, BookID, *Book],
) *userRepository {
	return &userRepository{
		m: map[UserID]*User{
			1: {
				ID:    1,
				Name:  "Alice",
				Books: booksLoadable,
			},
			2: {
				ID:    2,
				Name:  "Bob",
				Books: booksLoadable,
			},
		},
	}
}

var _ UserRepository = &userRepository{}

type userRepository struct {
	m map[UserID]*User
}

func (u *userRepository) List(ctx context.Context, ids []UserID) ([]*User, error) {
	users := make([]*User, len(ids))
	for i, id := range ids {
		users[i] = &User{
			ID:   id,
			Name: fmt.Sprintf("User %d", id),
		}
	}
	return users, nil
}

func (u *userRepository) All() ([]*User, error) {
	users := make([]*User, 0, len(u.m))
	for _, user := range u.m {
		users = append(users, user)
	}

	return users, nil
}

func NewBookRepository(
	placeLoadable preloader.HasOneLoadable[BookID, *Book, PlaceID, *Place],
) *bookRepository {
	return &bookRepository{
		m: map[BookID][]*Book{
			1: {
				&Book{
					ID:       1,
					Title:    "テスト本1",
					AuthorID: 1,
					PlaceID:  1,
					Place:    placeLoadable,
				},
				&Book{
					ID:       2,
					Title:    "テスト本2",
					AuthorID: 1,
					PlaceID:  2,
					Place:    placeLoadable,
				},
			},
		},
		authorLoadable: preloader.EmptyHasOneLoadable[BookID, *Book, UserID, *User](),
	}
}

type bookRepository struct {
	m              map[BookID][]*Book
	authorLoadable preloader.HasOneLoadable[BookID, *Book, UserID, *User]
}

var _ BookRepository = &bookRepository{}

func (b *bookRepository) List(ctx context.Context, ids []BookID) ([]*Book, error) {
	books := make([]*Book, 0, len(ids))
	for _, id := range ids {
		books = append(books, b.m[id]...)
	}

	return books, nil
}

func (b *bookRepository) ByUsers(ctx context.Context, userIDs []UserID) (map[UserID][]BookID, error) {
	m := map[UserID][]BookID{
		1: {
			1,
			2,
		},
	}

	res := make(map[UserID][]BookID, len(m))
	for _, userID := range userIDs {
		if bookIDs, ok := m[userID]; ok {
			res[userID] = bookIDs
		}
	}

	return res, nil
}

func (b *bookRepository) InjectAuthorLoadable(authorLoadable preloader.HasOneLoadable[BookID, *Book, UserID, *User]) {
	for _, book := range b.m {
		for _, b := range book {
			b.Author = authorLoadable
		}
	}
}

func NewPlaceRepository() *placeRepository {
	return &placeRepository{
		m: map[PlaceID]*Place{
			1: {
				ID:   1,
				Name: "Tokyo",
			},
			2: {
				ID:   2,
				Name: "Osaka",
			},
		},
	}
}

type placeRepository struct {
	m map[PlaceID]*Place
}

var _ PlaceRepository = &placeRepository{}

func (p *placeRepository) List(ctx context.Context, ids []PlaceID) ([]*Place, error) {
	places := make([]*Place, len(ids))
	for i, id := range ids {
		places[i] = p.m[id]
	}
	return places, nil
}
