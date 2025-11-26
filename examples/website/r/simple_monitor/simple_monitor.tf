resource "sakuracloud_simple_monitor" "foobar" {
  target = "www.example.com"

  delay_loop = 60
  timeout    = 10

  max_check_attempts = 3
  retry_interval     = 10

  health_check {
    protocol        = "https"
    port            = 443
    path            = "/"
    contains_string = "ok"
    status          = "200"
    host_header     = "example.com"
    sni             = true
    verify_sni      = true
    http2           = true
    # username    = "username"
    # password    = "password"
    # ftps        = "explicit"
  }

  description = "description"
  tags        = ["tag1", "tag2"]

  notify_email_enabled = true
  notify_email_html    = true
  notify_slack_enabled = true
  notify_slack_webhook = "https://hooks.slack.com/services/xxx/xxx/xxx"

  monitoring_suite {
    enabled = true
  }
}