# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [1.2.3] - 2022-11-07
### Fixed
- Fixed potential incorrect data when listing more than 100 records (reference overwrite)
- Fixed potential `nil pointer dereference` in versionService

## [1.2.2] - 2022-10-04
### Fixed
- Issue #15 : List functions are no longer limited to 100 results

## [1.2.1] - 2022-09-12
### Fixed
- Possible null pointer error in ClientApplication::GetSecret*

## [1.2.0] - 2022-07-21
### Added
- `Systems` Service for linked systems and groups on linked system
- `ClientApplications` Service to manage ClientApplications
- Read ClientApplications details
- Create new ClientApplication
- Read LinkedSystem details
- Read GroupOnSystem details
- Create new GroupOnSystem with optional `provgroups`
- Create Group with multiple admins and permissions for ClientApplications
