# withdrawal_service(出金サービス)

## 必要なもの

Docker, Docker Compose, Go 1.18

## システム概要

ユーザーが出金履歴登録リクエストをするとその履歴が`histories`テーブルに保存されるだけのシンプルなシステムです。ただし、合計 10 万円しか出金できないという制限があるので、10 万円を超えて登録しようとすると、登録はせずにステータス 400 を返すことを期待します。登録成功時に返却されるステータスコードは 201 です。

![sample drawio](https://user-images.githubusercontent.com/69618840/176994107-f3efb774-23a5-4930-86c2-019178897f69.png)

## 環境構築・テスト実行方法

`docker compose up`コマンドで DB サーバーと API サーバーを起動することができます。

API サーバーの起動が完了したら(`[GIN-debug] Listening and serving HTTP on :3000`のログ出力がされれば起動に成功しています。)

`go test ./...`コマンドによりテストを実行できます。テストを n 回連続で実行したい場合は`go test ./... -count n`のように実行してください。

API 実装はプルリクエスト単位で実装してあるのでテストを実行したいブランチに切り替えてからテストを実行してください。

[cosmtrek/air](https://github.com/cosmtrek/air)でホットリロードができるので、ブランチ変更後のビルドは不要です。

https://github.com/shuntagami/withdrawal_service/pull/1 の実装のテストを実行する例

```
$ gh pr checkout 1 # 作業ブランチに移動
$ go test ./... -count 10 # ホットリロードが完了してからテストを実行すること
```

## テストコード(main_test.go)の説明

user1 と user2 が出金履歴登録リクエストを同時に行うケースを想定しています。それぞれのユーザーは 10000 円の出金履歴登録リクエストを 12 回行うので、10 回のリクエストが成功し、残りの 2 回はステータスコード 400 で失敗となることを期待します。テストが fail するのは以下の 2 パターンです。

- (ユーザーごとに)10 万円を超える出金履歴が登録されていた場合

```
main_test.go:84:
        	Error Trace:	/Users/shuntagami/projects/withdrawal_service/main_test.go:84
        	Error:      	Should be true
        	Test:       	TestCreateHistory
        	Messages:   	user:1 amount 110000 over the amountLimit 100000

```

- API サーバーの返却するステータスコードが 201 もしくは 400 以外の場合

```
main_test.go:58:
        	Error Trace:	/Users/shuntagami/projects/withdrawal_service/main_test.go:58
        	            				/Users/shuntagami/projects/withdrawal_service/asm_arm64.s:1263
        	Error:      	Should be true
        	Test:       	TestCreateHistory
        	Messages:   	unexpected status code 500
```
