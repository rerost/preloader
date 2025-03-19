## computed_model
ActiveRecordの`users.preload(books: [:place, :author]` のようなことをGoでもやりたい。

ただし、GoではRepository層でモデルの解決がされることが多いため、こういったPreload実装が困難。
このライブラリでは、DBに依存せずRepositoryが実装しているすべてをPreloadできるようにするためのもの。

基本的なアイディア
モデルAからBへの依存関係があるとき、モデルAが存在する状態でBを取得するためには、
* モデルAからモデルBのIDへの変換 `[]A -> []BID`
* モデルBのIDからモデルBの解決 `[]BID -> []B`

この2つがあれば、実装できるはず。

参考: https://github.com/wantedly/computed_model

## TODO
- [x] `Preload(users, "Books", "Books.Place", "Books.Author")` をどうするか
- [x] HasOneの場合、使い勝手が悪い
- [x] LoadableのInjectが結構面倒。構造体に渡したり、相互参照があるときに面倒になってくる

## Example
```bash
cd example
go run main.go repository.go resource.go
```

## Type-Safe API

preloaderライブラリは、コンパイル時にLoadableの登録を強制する型安全なAPIを提供するようになりました。これにより、Loadableが登録されていない場合にコンパイルエラーが発生します。

### 例

```go
// 型安全なプロバイダーを作成
provider := preloader.NewTypedLoadableProvider()

// Loadableを登録
bookLoadable := preloader.EmptyLoadable[UserID, *TypedUser, BookID, *TypedBook]()
registeredBookLoadable := preloader.RegisterLoadable(provider, preloader.LoadableKey("Books"), bookLoadable)

// ユーザーを作成し、Loadableを登録
user := &TypedUser{ID: 1, Name: "Test User"}
user.SetProvider(provider)
user.RegisterBooksLoadable(registeredBookLoadable)

// "Books" Loadableが登録されているため、これはコンパイルされます
books, err := user.Books(context.Background())

// "Authors" Loadableが登録されていない場合、これはコンパイルエラーになります
book := &TypedBook{ID: 1, Title: "Test Book"}
book.SetProvider(provider)
// 以下の行はコンパイルエラーになります
// author, err := book.Author(context.Background())
```

### 型安全なモデルの定義

```go
// TypedUser は型安全なユーザーモデルです
type TypedUser struct {
    ID   UserID
    Name string
    
    provider *preloader.TypedLoadableProvider
    booksLoadable preloader.RegisteredLoadable[preloader.Registered, UserID, *TypedUser, BookID, *TypedBook]
}

// RegisterBooksLoadable は Books loadable を登録します
func (u *TypedUser) RegisterBooksLoadable(
    loadable preloader.RegisteredLoadable[preloader.Registered, UserID, *TypedUser, BookID, *TypedBook],
) {
    u.booksLoadable = loadable
}

// Books はユーザーの本を返します
func (u *TypedUser) Books(ctx context.Context) ([]*TypedBook, error) {
    loadable := preloader.GetRegisteredLoadable(u.provider, u.booksLoadable)
    return loadable.Load(ctx, u)
}
```

### ファントム型システム

新しいAPIは、ファントム型を使用してLoadableの登録状態をコンパイル時に検証します：

```go
// 登録状態を表すファントム型
type Registered struct{}
type NotRegistered struct{}

// 登録されたLoadableを表す型
type RegisteredLoadable[R any, ParentID comparable, Parent Resource[ParentID], NodeID comparable, Node Resource[NodeID]] struct {
    Key      LoadableKey
    Loadable Loadable[ParentID, Parent, NodeID, Node]
}
```

### 利点

- Loadable登録のコンパイル時型チェック
- 未登録のLoadableによる実行時エラーの排除
- 正しい使用法を強制する型安全なAPI
- ジェネリクスを活用した型安全なAPI設計
- Go 1.21との互換性
