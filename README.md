# gowiser
- [検索エンジン自作入門](https://gihyo.jp/book/2014/978-4-7741-6753-4)を参考に実装


## DB起動
```
docker compose up
```

## indexの作成
```
go run ./ i
```

## 単語検索
```
go run ./ s --keyword "Test body"
```
