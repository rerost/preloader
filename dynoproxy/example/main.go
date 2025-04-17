package main

import (
	"context"
	"fmt"
	"os"

	"github.com/rerost/preloader"
	"github.com/rerost/preloader/dynoproxy"
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
	
	placeLoadable, err := dynoproxy.DynamicHasOneLoadable(
		"Places",
		BookToPlace,
		placeRepository.List,
		true,
	)
	if err != nil {
		return err
	}
	
	bookRepository := NewBookRepository()
	bookLoader := &UsersToBooksLoader{bookRepository: bookRepository}
	bookLoadable, err := dynoproxy.DynamicLoadable(
		"Books",
		bookLoader.IDs,
		bookRepository.List,
	)
	if err != nil {
		return err
	}
	
	userRepo := NewUserRepository()
	authorLoadable, err := dynoproxy.DynamicHasOneLoadable(
		"Authors",
		BookToAuthor,
		userRepo.List,
		true,
	)
	if err != nil {
		return err
	}
	
	userBuilder := dynoproxy.NewModelBuilder().
		AddLoadable("Books", bookLoadable)
	
	bookBuilder := dynoproxy.NewModelBuilder().
		AddLoadable("Place", placeLoadable).
		AddLoadable("Author", authorLoadable)
	
	if err := userBuilder.Build(userRepo); err != nil {
		return err
	}
	if err := bookBuilder.Build(bookRepository); err != nil {
		return err
	}
	
	users, _ := userRepo.All()
	
	preloader := dynoproxy.NewPreloadHelper(ctx).
		AddLoadable(bookLoadable.Child(
			authorLoadable,
			placeLoadable,
		))
	
	if err := preloader.Preload(users); err != nil {
		return err
	}
	
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
	bookRepository *BookRepository
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

type UserID int

func (u *UserID) String() string {
	return fmt.Sprintf("User %d", *u)
}

type User struct {
	ID    UserID
	Name  string
	Books preloader.Loadable[UserID, *User, BookID, *Book]
}

func (u *User) GetResourceID() UserID {
	return u.ID
}

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

type UserRepository struct {
	m     map[UserID]*User
	Books preloader.Loadable[UserID, *User, BookID, *Book]
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		m: map[UserID]*User{
			1: {
				ID:   1,
				Name: "Alice",
			},
			2: {
				ID:   2,
				Name: "Bob",
			},
		},
	}
}

func (u *UserRepository) List(ctx context.Context, ids []UserID) ([]*User, error) {
	users := make([]*User, len(ids))
	for i, id := range ids {
		users[i] = &User{
			ID:    id,
			Name:  fmt.Sprintf("User %d", id),
			Books: u.Books,
		}
	}
	return users, nil
}

func (u *UserRepository) All() ([]*User, error) {
	users := make([]*User, 0, len(u.m))
	for _, user := range u.m {
		user.Books = u.Books
		users = append(users, user)
	}

	return users, nil
}

type BookRepository struct {
	m      map[BookID][]*Book
	Place  preloader.HasOneLoadable[BookID, *Book, PlaceID, *Place]
	Author preloader.HasOneLoadable[BookID, *Book, UserID, *User]
}

func NewBookRepository() *BookRepository {
	return &BookRepository{
		m: map[BookID][]*Book{
			1: {
				&Book{
					ID:       1,
					Title:    "テスト本1",
					AuthorID: 1,
					PlaceID:  1,
				},
				&Book{
					ID:       2,
					Title:    "テスト本2",
					AuthorID: 1,
					PlaceID:  2,
				},
			},
		},
	}
}

func (b *BookRepository) List(ctx context.Context, ids []BookID) ([]*Book, error) {
	books := make([]*Book, 0, len(ids))
	for _, id := range ids {
		for _, book := range b.m[id] {
			book.Place = b.Place
			book.Author = b.Author
			books = append(books, book)
		}
	}

	return books, nil
}

func (b *BookRepository) ByUsers(ctx context.Context, userIDs []UserID) (map[UserID][]BookID, error) {
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

type PlaceRepository struct {
	m map[PlaceID]*Place
}

func NewPlaceRepository() *PlaceRepository {
	return &PlaceRepository{
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

func (p *PlaceRepository) List(ctx context.Context, ids []PlaceID) ([]*Place, error) {
	places := make([]*Place, len(ids))
	for i, id := range ids {
		places[i] = p.m[id]
	}
	return places, nil
}
