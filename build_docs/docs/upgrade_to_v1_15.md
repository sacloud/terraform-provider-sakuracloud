# v1.15での変更点

---

## エンハンスドロードバランサでのレスポンスヘッダ設定

エンハンスドロードバランサ(sakuracloud_proxylb)にてポートごとに任意のレスポンスヘッダを付与できるようになりました。

#### 利用例

```hcl
resource "sakuracloud_proxylb" "example" {
  name           = "terraform-test-proxylb"
  plan           = 100
  vip_failover   = true
  sticky_session = true
  
  health_check {
    protocol    = "http"
    delay_loop  = 10
    host_header = "example.jp"
    path        = "/"
  }
  bind_ports {
    proxy_mode = "http"
    port       = 80
    
    # レスポンスヘッダ(複数指定可能:10個まで)
    response_header {
        header = "Cache-Control"
        value  = "public, max-age=10"
    }
  }
  servers {
      ipaddress = "${sakuracloud_server.server01.ipaddress}"
      port = 80
  }
}
```

