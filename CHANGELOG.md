# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]
## Added
- Add Basic Authentication support for project and issue CRUD requests (#3)
- Add option to skip SSL certificate verification (#4)

## [v0.1.0] - 2021-03-05
### Added
- Add direct project fields (#1)
   - Homepage
   - IsPublic
   - InheritMembers

### Removed
- Remove indirect project field CustomFields (#1)

### Fixed
- Project and Issue create, update, delete accept now additional positive HTTP codes HTTP 201 and HTTP 204