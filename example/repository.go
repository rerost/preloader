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
	provider preloader.LoadableProvider,
) *userRepository {
	return &userRepository{
		provider: provider,
		m: map[UserID]*User{
			1: {
				ID:       1,
				Name:     "Alice",
				provider: provider,
			},
			2: {
				ID:       2,
				Name:     "Bob",
				provider: provider,
			},
		},
	}
}

var _ UserRepository = &userRepository{}

type userRepository struct {
	m        map[UserID]*User
	provider preloader.LoadableProvider
}

func (u *userRepository) List(ctx context.Context, ids []UserID) ([]*User, error) {
	users := make([]*User, len(ids))
	for i, id := range ids {
		users[i] = &User{
			ID:       id,
			Name:     fmt.Sprintf("User %d", id),
			provider: u.provider,
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
	provider preloader.LoadableProvider,
) *bookRepository {
	return &bookRepository{
		provider: provider,
		m: map[BookID][]*Book{
			1: {
				&Book{
					ID:       1,
					Title:    "テスト本1",
					AuthorID: 1,
					PlaceID:  1,
					provider: provider,
				},
				&Book{
					ID:       2,
					Title:    "テスト本2",
					AuthorID: 1,
					PlaceID:  2,
					provider: provider,
				},
			},
		},
	}
}

type bookRepository struct {
	m        map[BookID][]*Book
	provider preloader.LoadableProvider
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

// InjectAuthorLoadable method is no longer needed with the LoadableProvider pattern

func NewPlaceRepository(
	provider preloader.LoadableProvider,
) *placeRepository {
	return &placeRepository{
		provider: provider,
		m: map[PlaceID]*Place{
			1: {
				ID:       1,
				Name:     "Tokyo",
				provider: provider,
			},
			2: {
				ID:       2,
				Name:     "Osaka",
				provider: provider,
			},
		},
	}
}

type placeRepository struct {
	m        map[PlaceID]*Place
	provider preloader.LoadableProvider
}

var _ PlaceRepository = &placeRepository{}

func (p *placeRepository) List(ctx context.Context, ids []PlaceID) ([]*Place, error) {
	places := make([]*Place, len(ids))
	for i, id := range ids {
		places[i] = &Place{
			ID:       id,
			Name:     p.m[id].Name,
			provider: p.provider,
		}
	}
	return places, nil
}
