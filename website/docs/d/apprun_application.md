---
layout: "sakuracloud"
page_title: "SakuraCloud: sakuracloud_apprun_application"
subcategory: "AppRun"
description: |-
  Get information about an existing AppRun Application.
---

# Data Source: sakuracloud_apprun_application

Get information about an existing AppRun Application.

## Argument Reference

* `name` - (Required) The name of application.

## Attribute Reference

* `id` - The id of the AppRun Application.
* `timeout_seconds` - The time limit between accessing the application's public URL, starting the instance, and receiving a response.
* `port` - The port number where the application listens for requests.
* `min_scale` - The minimum number of scales for the entire application.
* `max_scale` - The maximum number of scales for the entire application.
* `components` - The application component information.
* `public_url` - The public URL.
* `status` - The application status.


---

A `components` block exports the following:

* `name` - The component name.
* `max_cpu` - The maximum number of CPUs for a component.
* `max_memory` - The maximum memory of component.
* `deploy_source` - The sources that make up the component.
* `env` - The environment variables passed to components.
* `probe` - The component probe settings.

---

A `deploy_source` block supports the following:

* `container_registry` - A `container_registry` block as defined below.

---

A `container_registry` block exports the following:

* `image` - The container image name.
* `server` - The container registry server name.
* `username` - The container registry credentials.

---

A `env` block supports the following:

* `key` - The environment variable name.
* `value` - environment variable value.

---

A `probe` block exports the following:

* `http_get` - A list of `http_get` blocks as defined below.

---

A `http_get` block exports the following:

* `path` - The path to access HTTP server to check probes.
* `port` - The port number for accessing HTTP server and checking probes.
* `headers` - One or more `headers` blocks as defined below.

---

A `headers` block supports the following:

* `name` - The header field name.
* `value` - The header field value.
