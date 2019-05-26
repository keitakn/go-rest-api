# go-rest-api
GoLang製のREST APIサンプルプロジェクト

## ローカル環境で動作させる

### 初回起動時

`docker-compose up --build -d`

### 2回目以降の起動時

`docker-compose up`

`curl -v http://127.0.0.1:9999` で動作している事を確認出来る。

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
