# go-rest-api
GoLang製のREST APIサンプルプロジェクト

## ローカル環境で動作させる

### 初回起動時

`docker-compose up --build -d`

### 2回目以降の起動時

`docker-compose up`

`curl -v http://127.0.0.1:8080` で動作している事を確認出来る。

### 停止

`docker-compose down`

### 停止とDockerイメージの削除を同時に行う

`docker-compose down --rmi all`

## 本番環境で動作させる

以下のコマンドを実行しコンテナを立ち上げる。

```
docker build -t go-rest-api .
docker run -p 8080:8080 -d --name go-rest-api go-rest-api
```

`curl -v http://127.0.0.1:8080` で動作確認を行い問題なければDocker HubやECR等に反映させる。

## コードの整形
`go fmt`

## ECRへのプッシュ

`docker-push-ecr.sh` を実行して下さい。（当然事前にECRリポジトリを作成しておく必要があります）

このscriptの実行には環境変数 `AWS_ACCOUNT_ID` を設定しておく必要があります。

[direnv](https://github.com/direnv/direnv) 等を使って環境変数を設定して下さい。

## 各エンドポイントのメモ

## ユーザー登録

```
curl -X POST -v \
-H "Content-type: application/json" \
-d \
'
{
  "name": "keitakn"
}
' \
http://127.0.0.1:8080/users
```

## ユーザー取得（すべて）

`curl -v http://127.0.0.1:8080/users`

## ユーザー取得（ID指定）

`curl -v http://127.0.0.1:8080/users/{id}`


## ユーザー更新

```
curl -X PUT -v \
-H "Content-type: application/json" \
-d \
'
{
  "name": "keita-keita"
}
' \
http://127.0.0.1:8080/users/{id}
```

## ユーザー削除

`curl -v -X DELETE http://127.0.0.1:8080/users/2`
