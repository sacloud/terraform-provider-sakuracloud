# v1.1での変更点

---

## シンプル監視/GSLBでのヘルスチェックのデータ型変更

v1.1ではシンプル監視/GSLBでのヘルスチェック(`health_check`)についてデータ型の変更が行われています。  

tfファイルの変更は不要ですが、`terraform.tfstate`の内容が変更されているために
`plan`や`apply`で変更が必要と表示されます。

### 対応方法

`terraform apply`を実行しリソースの更新を行ってください。

## サーバのNIC(1番目のNIC)を共有セグメント/スイッチに接続しない方法の変更

これまでは以下のように`nic`に空文字を指定することで何処とも接続されていないNICを追加可能でした。

```hcl
resource sakuracloud_server "server" {
  nic = ""
}
```

v1.1以降は以下のように`disconnect`と指定する必要があります。

```hcl
resource sakuracloud_server "server" {
  nic = "disconnect"
}
```

