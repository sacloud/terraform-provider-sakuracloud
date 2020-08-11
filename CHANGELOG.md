## 2.3.6 (Unreleased)

FIXES

* Fix error handling at VPCRouter [GH-757] (@yamamoto-febc)
* Use libsacloud v2.6.4 [GH-754] (@yamamoto-febc)
* Fix time comparison method - use time#Equal() [GH-753] (@yamamoto-febc)
* Use libsacloud v2.6.3  [GH-751] (@yamamoto-febc)

DOCS

* docs: Terraform v0.13 [GH-758] (@yamamoto-febc)

## 2.3.5 (2020-06-19)

* Update dependencies - libsacloud v2.6.1 [GH-748] (@yamamoto-febc)

## 2.3.4 (2020-06-16)

* ProxyLB: supports anycast [GH-747] (@yamamoto-febc)

## 2.3.3 (2020-05-18)

Note: Publishing in the Terraform Registry is supported this version and later.
see http://registry.terraform.io/providers/sacloud/sakuracloud/

FEATURES

* Support for publishing in the Terraform Registry [GH-744] (@yamamoto-febc)

FIXES

* Fixes plan changing of ProxyLB [GH-745] (@yamamoto-febc)

## 2.3.2 (2020-05-15)

* This is an experimental release for testing publishing to the Terraform registry. Don't use this in a production environment.

## 2.3.1 (2020-04-24)

FIXES

* Modify how to determine whether to pass disk_edit_parameter to ServerBuilder [GH-737] (@yamamoto-febc)

MISC

* Fix broken CI - install golangci-lint via install script [GH-735] (@yamamoto-febc)

## 2.3.0 (2020-04-20)

* Startup Script Parameters [GH-731] (@yamamoto-febc)
    * libsacloud v2.5.1
* libsacloud v2.5.2 - improve error messages [GH-733] (@yamamoto-febc)    

## 2.2.0 (2020-03-17)

FEATURES

* Add sakuracloud_archive_share resource [GH-728] (@yamamoto-febc)
* Supports transferred/shared archives [GH-727] (@yamamoto-febc)
    * libsacloud v2.4.1

IMPROVEMENTS

* Set ID to state even if got error from builders [GH-726] (@yamamoto-febc)
* libsacloud v2.3.0 - MariaDB 10.4 [GH-724]

## 2.1.2 (2020-03-10)

* Remove deletion waiter [GH-713] (@yamamoto-febc)
* libsacloud v2.1.7 [GH-713] (@yamamoto-febc)
* Go 1.14 [GH-712] (@yamamoto-febc)
* Fix wrong error message [GH-718] (@yamamoto-febc)
* libsacloud v2.1.8 to avoid marshal JSON error at SIM [GH-714] (@yamamoto-febc)
* libsacloud v2.1.9 [GH-723] (@yamamoto-febc)

## 2.1.1 (2020-02-28)

IMPROVEMENS/FIXES

* tfproviderlint v0.10.0 [GH-708] (@yamamoto-febc)
* libsacloud v2.1.4 [GH-708] (@yamamoto-febc)
* Upgrade libsacloud to v2.1.5 [GH-709] (@yamamoto-febc)

## 2.1.0 (2020-02-14)

FEATURES

* Container Registry: VirtualDomain/User permission [GH-704] (@yamamoto-febc)
* PostgreSQL 12.1 [GH-704] (@yamamoto-febc)

IMPROVEMENTS

* Terraform Plugin SDK v1.7.0 [GH-703] (@yamamoto-febc)
* tfproviderlint v0.9.0 [GH-698] (@yamamoto-febc)

## 2.0.1 (2020-02-06)

FIXES

* libsacloud v2.0.2 [GH-697] @yamamoto-febc

IMPROVEMENTS

* terraform-plugin-sdk v1.6.0 and tfproviderlint v0.9.0 [GH-698] @yamamoto-febc
* Use libsacloud v2.0.1 [GH-696] @yamamoto-febc

## 2.0.0 (2020-01-31)

NOTES:

* Initial release of v2.0

