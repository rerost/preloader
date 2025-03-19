package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rerost/preloader"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	placeRepository := NewPlaceRepository()
	placeLoadable := preloader.NewHasOneLoadable("Places", BookToPlace, placeRepository.List, true)

	bookRepository := NewBookRepository(placeLoadable)
	bookLoader := UsersToBooksLoader{bookRepository}
	bookLoadable := preloader.NewLoadable("Books", bookLoader.IDs, bookRepository.List)

	userRepo := NewUserRepository(bookLoadable)
	authorLoadable := preloader.NewHasOneLoadable("Authors", BookToAuthor, userRepo.List, true)
	bookRepository.InjectAuthorLoadable(authorLoadable)

	users, _ := userRepo.All()

	// Preload
	if err := preloader.Preload(
		ctx,
		users,
		bookLoadable.Child(
			authorLoadable,
			placeLoadable,
		),
	); err != nil {
		return err
	}

	// Print
	for _, user := range users {
		books, err := user.Books.Load(ctx, user)
		if err != nil {
			return err
		}
		for _, book := range books {
			place, err := book.Place.Load(ctx, book)
			if err != nil {
				return err
			}
			author, err := book.Author.Load(ctx, book)
			if err != nil {
				return err
			}
			fmt.Printf(
				"ユーザー名: %v, タイトル: %v, 場所ID: %v, 場所: %v, 著者ID: %v, 著者: %v\n",
				user.Name,
				book.Title,
				place.ID,
				place.Name,
				author.ID,
				author.Name,
			)
		}
	}
	return nil
}

type UsersToBooksLoader struct {
	bookRepository BookRepository
}

func (u *UsersToBooksLoader) IDs(ctx context.Context, users []*User) (map[UserID][]BookID, error) {
	userIDs := make([]UserID, len(users))
	for i, user := range users {
		userIDs[i] = user.ID
	}

	resMap, err := u.bookRepository.ByUsers(ctx, userIDs)
	if err != nil {
		return nil, err
	}

	return resMap, nil
}

func BookToPlace(ctx context.Context, books []*Book) (map[BookID][]PlaceID, error) {
	res := make(map[BookID][]PlaceID, len(books))

	for _, book := range books {
		res[book.ID] = []PlaceID{book.PlaceID}
	}

	return res, nil
}

func BookToAuthor(ctx context.Context, books []*Book) (map[BookID][]UserID, error) {
	res := make(map[BookID][]UserID, len(books))

	for _, book := range books {
		res[book.ID] = []UserID{book.AuthorID}
	}

	return res, nil
}
