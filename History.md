# 0.7.6 / 2026-08-24

> [!IMPORTANT]
> This release changes the filenames inside the archives. Previously you had `envrun_linux_amd64`, now each archive only contains the `envrun` binary (no more appended OS and ARCH).

  * chore: replace old build-system

# 0.7.5 / 2026-08-24

  * chore: update dependencies, replace go_helpers/v2 module
  * fix(deps): update module github.com/luzifer/go-openssl/v4 to v4.2.5 (#2)
  * fix(deps): update module github.com/luzifer/rconfig/v2 to v2.6.2 (#4)
  * fix(deps): update module github.com/sirupsen/logrus to v1.10.1 (#6)

# 0.7.4 / 2024-12-12

  * Update Go dependencies

# 0.7.3 / 2024-02-28

  * Update dependencies

# 0.7.2 / 2023-12-19

  * Update dependencies

# 0.7.1 / 2023-11-04

  * Update dependencies
  * Fix: `fmt.Errorf` behaves different than `errors.Wrap`

# 0.7.0 / 2023-09-30

  * Update dependencies, fix linter errors
  * Replace repo-runner CI

# 0.6.2 / 2020-02-13

  * Fix: Panic on invalid entries in env-file

# 0.6.1 / 2019-08-04

  * Fix: Add check for missing command arguments
  * Add go module support, update dependencies
  * Remove old vendoring

# 0.6.0 / 2018-11-02

  * Update dependencies
  * Update OpenSSL library, add CLI support for openssl-sha256

# 0.5.0 / 2018-09-14

  * Update dependencies
    * Adds OpenSSL >=1.1.0c support
  * Update README

# 0.4.2 / 2018-06-01

  * Fix: Report gpg-symmetric encryption in help

# 0.4.1 / 2018-06-01

  * Fix: Infinite loop on wrong password
  * Fix: All methods do not allow trailing spaces, remove them

# 0.4.0 / 2018-06-01

  * Add support for symmetric GPG encryption
  * Allow reading passphrase from file
  * Rework logging
  * Migrate to `dep` for vendoring

# 0.3.1 / 2017-03-21

  * Add Github publishing
  * Vendor dependencies

# 0.3.0 / 2017-03-21

  * Add support for encrypted env files


0.2.0 / 2016-04-15
==================

  * Allow comments in .env file
  * Fix: Missing link in README

0.1.0 / 2016-02-06
==================

  * Initial version
