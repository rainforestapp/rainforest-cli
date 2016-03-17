# Rainforest CLI Changelog

## 1.2.1 - 19th February 2016
- Fixed a bug where uploading was stuck in an infinite loop if an embedded id did not exist (7b02b2f66dbd47098a7c1d5f79bc60a0cbe8984f, @epaulet)
- Fixed a bug that occurred when specifying a nested test folder without creating parent folders first (6c1b0e02c858f9d9c264e771f964b3e1a4ea8c7e, @epaulet)
- Removed 'ro' tag and use 'rainforest-cli' as the test's source instead.
(10864a7e054d4c869f6a345608b2d1c1c0925fe8, @epaulet)
- Add retries on API calls so that minor interruptions do not cancel Rainforest builds.
(98021337a3fbbf16c7cd858bbec5d925fb86c939, @epaulet)

## 1.2.0 - 8th February 2016
- Add support for embedded tests.
- Add support for customizable RFML ids.

## 1.1.4 - 15th January 2016
- Customizable folder location for rainforest tests (fa4418738311cee8ca25cbb22a8ca52aa9cbd873, @ukd1)
- Update valid browser list, though this doesn't include custom browsers today (e6195c42f95cce72a17f49643bfe8c297baf8dd9, @ukd1)

## 1.1.2 - 15th January 2016
- Fixed specs (7c4af508d8cfa95363ee9976f1fa6f01f7c8d27b, @ukd1)
- Fixed a bug for users of Pivotal which caused tags to be incorrectly parsed (82966e7e739b590b396266d12d72605d6e19c12b, @chaselee)

## 1.1.0 - 6th January 2016
- Added support for tests stored in plain text format, so you can store them locally

## 1.0.5 - 25th August 2015
- Added environment support (278f4fe9a1ca9f507fe1e4b11935d9c37056786b)

## 0.0.11 - 30th Sept 2013
- First changelog entry.
- Fixed travis config (9f74cb37bc451356458d8087ebd2694271eedfc2)
