# practice-golang-api

[APIを作りながら進むGo中級者への道](https://techbookfest.org/product/jXDAEU1dR53kbZkgtDm9zx?productVariantID=dvjtgpjw8VDTXNqKaanTVi)を元にGoを使ってAPIサーバーを構築する練習。

## memo
### 12. ユーザー認証
OAuth同意画面の作成ではprofileだけはなく、openidもスコープに追加する。  
IDトークンの確認では以下のように`nonce`の値と`state`の値を追加する(`state`はなくても動作するかもしれないが`nonce`は必須)。  
```
https://accounts.google.com/o/oauth2/v2/auth?client_id=[クライアントID]&response_type=id_token&scope=openid%20profile&redirect_uri=http://localhost:8080/callback&nonce=random_nonce_value&&state=random_state_value
```
