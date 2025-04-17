# preloader-dynoproxy

`dynoproxy` は `go-dyno` を使用して `preloader` をより使いやすくするためのパッケージです。

## 特徴

- `go-dyno` の動的プロキシ機能を使用して、`Loadable` と `HasOneLoadable` インターフェースの作成を簡素化
- `ModelBuilder` を使用して、モデルへの Loadable の注入を自動化
- `PreloadHelper` を使用して、プリロード処理を簡素化
- リフレクションを使用して、構造体フィールドへの Loadable の自動注入をサポート

## 使用例

```go
// Create loadables using dynoproxy
placeLoadable, err := dynoproxy.DynamicHasOneLoadable(
    "Places",
    BookToPlace,
    placeRepository.List,
    true,
)

// Use ModelBuilder to inject loadables
bookBuilder := dynoproxy.NewModelBuilder().
    AddLoadable("Place", placeLoadable).
    AddLoadable("Author", authorLoadable)

// Inject loadables into repository
if err := bookBuilder.Build(&bookRepository); err != nil {
    return err
}

// Preload using PreloadHelper
preloader := dynoproxy.NewPreloadHelper(ctx).
    AddLoadable(bookLoadable.Child(
        authorLoadable,
        placeLoadable,
    ))

if err := preloader.Preload(users); err != nil {
    return err
}
```

## 従来の方法との比較

### 従来の方法

```go
placeLoadable := preloader.NewHasOneLoadable("Places", BookToPlace, placeRepository.List, true)

bookRepository := NewBookRepository(placeLoadable)
bookLoader := UsersToBooksLoader{bookRepository}
bookLoadable := preloader.NewLoadable("Books", bookLoader.IDs, bookRepository.List)

userRepo := NewUserRepository(bookLoadable)
authorLoadable := preloader.NewHasOneLoadable("Authors", BookToAuthor, userRepo.List, true)
bookRepository.InjectAuthorLoadable(authorLoadable)

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
```

### dynoproxy を使用した方法

```go
// Create loadables using dynoproxy
placeLoadable, err := dynoproxy.DynamicHasOneLoadable(
    "Places",
    BookToPlace,
    placeRepository.List,
    true,
)

bookRepository := NewBookRepository()
bookLoader := &UsersToBooksLoader{bookRepository: bookRepository}
bookLoadable, err := dynoproxy.DynamicLoadable(
    "Books",
    bookLoader.IDs,
    bookRepository.List,
)

userRepo := NewUserRepository()
authorLoadable, err := dynoproxy.DynamicHasOneLoadable(
    "Authors",
    BookToAuthor,
    userRepo.List,
    true,
)

// Use ModelBuilder to inject loadables
userBuilder := dynoproxy.NewModelBuilder().
    AddLoadable("Books", bookLoadable)

bookBuilder := dynoproxy.NewModelBuilder().
    AddLoadable("Place", placeLoadable).
    AddLoadable("Author", authorLoadable)

// Inject loadables into repositories
if err := userBuilder.Build(&userRepo); err != nil {
    return err
}
if err := bookBuilder.Build(&bookRepository); err != nil {
    return err
}

// Preload using PreloadHelper
preloader := dynoproxy.NewPreloadHelper(ctx).
    AddLoadable(bookLoadable.Child(
        authorLoadable,
        placeLoadable,
    ))

if err := preloader.Preload(users); err != nil {
    return err
}
```

## メリット

1. Loadable の作成と注入が簡素化され、コードの可読性が向上
2. 相互参照がある場合でも、ModelBuilder を使用して簡単に Loadable を注入可能
3. 構造体フィールドへの自動注入により、手動での注入コードが不要
4. PreloadHelper を使用して、プリロード処理を簡素化
