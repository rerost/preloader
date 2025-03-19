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
- [ ] LoadableのInjectが結構面倒。構造体に渡したり、相互参照があるときに面倒になってくる

## Example
```bash
cd example
go run main.go repository.go resource.go
```
