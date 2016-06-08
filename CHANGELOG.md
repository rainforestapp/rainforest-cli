# Rainforest CLI Changelog

## 1.6.0 - 8th June 2016
- Add `rm` command to delete tests. (cf077bd440e83852d8d23f66eb72ba94fece0c07, @marktran)
- New tests use title given when created using `rainforest new [TEST_NAME]`.
(01e30ba9d97558ba11c239a5c9842192d38dfd3f, @epaulet)
- Remove static browser list and stop client side validation of browsers for run
creation. Allow the API to validate against the dynamic list of client browsers
instead. (48a4d11e524d020e78e14991bf8a0c5bf82b65c9, @epaulet)

## 1.5.0 - 23rd May 2016
- Retry on API exceptions when creating a run. (85830adcef426e64bd72c9d6208881e955a5bb0c, @bbeck)
- Add `browsers` command. (2a810ec27edfc66ef7bf27d8cb7da1129b05e32b, @epaulet)
- Add support for files using `app-source-url`. (562e4772e71e8028209d128091ff644f4ae0a9f6, @marianafranco)
- Remove newlines from test actions and questions when exporting. (e28583b553b5f30b33b232b2e377c109123b11ff, @epaulet)

## 1.4.0 - 22nd April 2016
- Support for new `--version` command. (4362c85fe599a02eaa1b772d184be31e692e934e, @epaulet)
- Validate duplicate RFML IDs before uploading. (67f71d053c755eaf92c1bd205931e89e903b88c9, @curtis-rainforestqa)
- Add `folders` commands for a folder ID reference. (4ab19fec0924b4764e140fb3c5aa85f1dbfe4006, @epaulet)

## 1.3.1 - 11th April 2016
- Support crowd selection. (03fedacfb7a6e69a174fb3e0e1fd75218fdbbfa9, @ukd1)

## 1.3.0 - 7th April 2016
- Export tests with embedded tests unflattened. (0ed4c62cac8a0d5fbd98f03190d3c18c48ac7119,
@epaulet)
- Add option to save API token in environment rather than specifying in commands.
(eaa32e87dff2881074c920f6ffc278d1fcd25ae7, @valeriangalliat)
- Fixed a bug where tests would upload without the correct source attribute if
an upload command failed in the middle of executing. (92df14606304957c5c58719a8999471df5f4f8c0,
@epaulet)
- Specify app source url as a command line option. (d02f750e885824c1b6f141344af9a34fc99e7527,
@ukd1)
- Add support for redirect flag on steps and embedded test. (e54a3f78333d4b8398b8aece40ebfbaaf4113eb4,
@epaulet)
- Add support for site_id attribute in RFML test and add `sites` command for
site ID reference. (7b628b12879f5c2230181d5e4badf785c26c8035, @epaulet)
- Add support for exporting specified tests using test IDs found on dashboard.
(69d104d7452dcb2ba7925d1de86532f250b72f41, @ziahamza)

## 1.2.2 - 21st March 2016
- Add support for Ruby 1.9.3 for easier usage on CircleCI. (16d74306a160c0fca8d34bc32493119051179c90, @epaulet)

## 1.2.1 - 18th March 2016
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
