---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_apprun_application"
subcategory: "AppRun"
description: |-
  Manages a SakuraCloud AppRun Application.
---

# sakuracloud_apprun_application

Manages a SakuraCloud AppRun Application.

## Example Usage

```hcl
resource "sakuracloud_apprun_application" "foobar" {
  name            = "foobar"
  timeout_seconds = 60
  port            = 80
  min_scale       = 0
  max_scale       = 1
  components {
    name       = "foobar"
    max_cpu    = "0.1"
    max_memory = "256Mi"
    deploy_source {
      container_registry {
        image    = "foorbar.sakuracr.jp/foorbar:latest"
        server   = "foorbar.sakuracr.jp"
        username = "user"
        password = "password"
      }
    }
    env {
      key   = "key"
      value = "value"
    }
    env {
      key   = "key2"
      value = "value2"
    }
    env {
      key   = "key3"
      value = "value3"
    }
    probe {
      http_get {
        path = "/"
        port = 80
        headers {
          name  = "name"
          value = "value"
        }
        headers {
          name  = "name2"
          value = "value2"
        }
      }
    }
  }

  traffics {
    version_index = 0
    percent       = 100
  }
}
```

## Argument Reference

* `name` - (Required) The name of application.
* `timeout_seconds` - (Required) The time limit between accessing the application's public URL, starting the instance, and receiving a response.
* `port` - (Required) The port number where the application listens for requests.
* `min_scale` - (Required) The minimum number of scales for the entire application.
* `max_scale` - (Required) The maximum number of scales for the entire application.
* `components` - (Required) The application component information.
* `traffics` - (Optional) The application traffic.


---

A `components` block supports the following:

* `name` - (Required) The component name.
* `max_cpu` - (Required) The maximum number of CPUs for a component. The values in the list must be in [`0.1`/`0.2`/`0.3`/`0.4`/`0.5`/`0.6`/`0.7`/`0.8`/`0.9`/`1`].
* `max_memory` - (Required) The maximum memory of component. The values in the list must be in [`256Mi`/`512Mi`/`1Gi`/`2Gi`].
* `deploy_source` - (Required) The sources that make up the component.
* `env` - (Optional) The environment variables passed to components.
* `probe` - (Optional) The component probe settings.

---

A `deploy_source` block supports the following:

* `container_registry` - (Optional) A `container_registry` block as defined below.

---

A `container_registry` block supports the following:

* `image` - (Required) The container image name.
* `server` - (Optional) The container registry server name.
* `username` - (Optional) The container registry credentials.
* `password` - (Optional) The container registry credentials.

---

A `env` block supports the following:

* `key` - (Optional) The environment variable name.
* `value` - (Optional) environment variable value.

---

A `probe` block supports the following:

* `http_get` - (Required) A `http_get` block as defined below.

---

A `http_get` block supports the following:

* `path` - (Required) The path to access HTTP server to check probes.
* `port` - (Required) The port number for accessing HTTP server and checking probes.
* `headers` - (Optional) One or more `headers` blocks as defined below.

---

A `headers` block supports the following:

* `name` - (Optional) The header field name.
* `value` - (Optional) The header field value.

---

A `traffics` block supports the following:

~> **Note:** When an application is created or updated, its configuration information is stored as a version. version_index specifies the index of the list of versions, sorted in descending order by creation date. For example, if there are three versions, "version_index = 0" refers to the most recent version, and "version_index = 2" refers to the oldest version.

* `version_index` - (Required) The application version index.
* `percent` - (Required) The percentage of traffic dispersion.


### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration/resources.html#operation-timeouts) for certain actions:

* `create` - (Defaults to 5 minutes) Used when creating the AppRun Application


* `update` - (Defaults to 5 minutes) Used when updating the AppRun Application

* `delete` - (Defaults to 20 minutes) Used when deleting AppRun Application



## Attribute Reference

* `id` - The id of the AppRun Application.
