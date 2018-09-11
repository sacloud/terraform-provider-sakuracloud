# v1.4での変更点

---

## marker_tags機能の除去

v1.4ではmarker_tags機能が除去されました。  
これは、プロバイダー設定ブロック内で指定しておくと各リソースに対し一律で任意のタグを付与できると言う機能でした。  
主にTerraformで作成されたリソースであることをタグで示したい場合などに利用されていました。

```hcl
provider sakuracloud {
  use_marker_tags = true
  marker_tag_name = "@terraform"
}
```

### 対応方法

各リソース側で`tags`を指定するようにしてください。

