# v1.8での変更点

---

## Terraform v0.12対応

このバージョンからTerraform v0.12に対応しています。  
Terraform v0.12ではtfファイル記法などに大幅な変更が行われました。  
変更内容の詳細については[公式ドキュメント](https://terraform.io/docs/)を参照してください。

このドキュメントに記載しているtfファイルの例はTerraform v0.12対応となっています。

Terraform v0.12以前のバージョンをご利用中の場合、Terraform v0.12を用いて以下のコマンドを実行することでtfファイルのマイグレーションが行えます。

```bash
# terraform v0.12以前のバージョン向けのtfファイルをv0.12向けに書き換え
$ terraform 0.12upgrade
```
