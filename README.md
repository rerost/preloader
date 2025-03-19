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

preloaderライブラリは、コンパイル時にLoadableの登録を強制する型安全なAPIを提供するようになりました。これにより、Loadableが登録されていない場合に実行時エラーが発生するのを防ぎます。

### 例

```go
// 型安全なプロバイダーを作成
provider := preloader.NewTypedLoadableProvider()

// Loadableを登録
bookLoadable := preloader.NewLoadable("Books", bookLoader.IDs, bookRepository.List)
provider.RegisterTypedLoadable(preloader.LoadableKey("Books"), bookLoadable)

// "Books" Loadableが登録されているため、これはコンパイルされます
user := &TypedUser{ID: 1, Name: "Test User"}
user.SetProvider(provider)
books, err := user.Books(context.Background())

// "Authors" Loadableが登録されていない場合、これはコンパイルエラーになります
book := &TypedBook{ID: 1, Title: "Test Book"}
book.SetProvider(provider)
author, err := book.Author(context.Background())
```

### 利点

- Loadable登録のコンパイル時型チェック
- 未登録のLoadableによる実行時エラーの排除
- 正しい使用法を強制する型安全なAPI
