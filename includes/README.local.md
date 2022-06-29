## terraform-provider-sakuracloud独自の項目

ゴール上書きの警告を防ぐため、go/common.mkの`test`と`testacc`を以下のように修正しています。

```makefile
.PHONY: test
test:
	TF_ACC= SAKURACLOUD_APPEND_USER_AGENT="$(UNIT_TEST_UA)" go test -v $(TESTARGS) -timeout=30s ./...

.PHONY: testacc
testacc:
	TF_ACC=1 SAKURACLOUD_APPEND_USER_AGENT="$(ACC_TEST_UA)" go test -v $(TESTARGS) -timeout 240m ./...

```

今後sacloud/makefileの更新を取り入れる際に注意が必要です。