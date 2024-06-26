# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]


## [1.3.5] - 2024-06-25
### Changed
- Issue # : Vaults and groups always uses latest as contract id

## [1.3.4] - 2024-06-25
### Fixed
- Issue # : Fix incorrect check for no supported versions

## [1.3.3] - 2024-06-25
### Changed
- Issue # : Optional services based on supported version of keyhub

## [1.3.2] - 2024-06-12
### Changed
- Issue # : Add deprecation warning to readme
- Issue # : retract version 1.3.1 due to contract version/model mismatch
- Issue # : Switch group and vault service to contract version 71
- Issue # : Downgrade other services to contract version 60

## [1.3.1] - 2023-05-09
### Changed
- Issue # : Build with 1.21
- Issue #27 : Missing contract headers
- Issue #17 : implement service accounts 
- Issue #14 : create a launchpadtile


## [1.3.0] - 2023-05-09
### Changed
- Issue # : Require and changes for keyhub contract version 62 to fix problem with provision groups
### Added
- Issue #24 : Allow creation of GroupOnSystem without provGroups

## [1.2.5] - 2023-03-10
### Fixed
- Hotfix: Fix undetected parse error for vault records without an enddate

### Added
- Issue #17 : Implement Service Accounts

## [1.2.4] - 2023-02-23
### Changed
- Issue #18 : Improved error reporting by returning the ErrorReport
### Fixed
- Issue #20 : Fixed index-out-of-range panic while getting vaultRecord  
- Issue #21 : Can't set expire / end date for vault records

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
