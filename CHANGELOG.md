# Changelog
All notable changes to this project will be documented in this file
	see also https://github.com/clojure/brew-install/blob/1.10.1/CHANGELOG.md for
changes introduced by the official installer
## next

## [1.10.1.507] 2020-02-08
### Changed
- Update to tools.deps.alpha 0.8.661
- Add -Sthreads option for concurrent downloads
- Use -- to separate dep options and clojure.main options
- Report clj version in clj -h
	
### Changed
Remove -Xms on tool jvms

## [1.10.1.489] 2019-11-27
### Changed
- Added -Strace option
- Update to tools.deps.alpha 0.8.599
## [1.10.1.478] - 2019-10-19
### Changed
- Adapt to version 1.10.1.478 of the official installer

## [1.10.1.469] - 2019-08-12

### Changed
- Update to tools.deps.alpha 0.7.541
- Add slf4j-nop which was removed from tools.deps.alpha
- Use min heap setting, not max, on internal tool calls
