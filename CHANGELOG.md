# Rainforest CLI Changelog

## 1.12.3 - 3rd March 2017
- Validate RFML tests for a title.
(beb439a3d2c12f618d6b4b00a3a08f1e37bbbe7a, @epaulet)

## 1.12.2 - 23rd February 2017
- Default `start_uri` attribute to "/" if omitted from RFML test.
(ec6049407ec635b7f7bc4f8da516ccddad2b78b3, @epaulet)

## 1.12.1 - 7th February 2017
- Check columns returned from API before uploading rows when creating generators.
(458fdf5bf3ba28588dfb48d1253192ba477ac7ae, @epaulet)

## 1.12.0 - 6th February 2017
- Add category column to list of sites when using the `sites` command.
(11c24c3484a587b7305dcf279872fa5e5be02b5a, @epaulet)

## 1.11.0 - 4th November 2016
- Add `--single-use` flag for CSV uploads. (17a19694a788365beb59e634bd7286c86528484b, @epaulet)

## 1.10.4 - 31st October 2016
- Do not upload test metadata as test description. (1c8b46d2156da17be939ac67edefe2bd2b56af47, @epaulet)

## 1.10.3 - 21st October 2016
- Fix bug with accidental removal of 3 commands. (d9ebbd702ba1b1b64b9a1a222c301315b89d743e, @curtis-rainforestqa)
- Add better tolerance for server errors. (b4180332ca13fec0972c3d65dbce00242907ca38, @epaulet)

## 1.10.2 - 11th October 2016
- Limit CSV batch upload size to chunks of 20 at a time for more reliability.
(abfed07df5635eb88029ee0b6cf8eea3a538fff6, @epaulet)

## 1.10.1 - 10th October 2016
- Add documentation for commands to `--help` text. (fe3fddbe7086a6b914ffe74f65ae128a14277f82, @epaulet)
- Use more efficient CSV upload method. (e716e76b7627b334c4f8b50fd0afdad6fadd162e, @epaulet)

## 1.10.0 - 19th September 2016
- Add `--overwrite-variable` option to overwrite existing tabular variables if the
desired variable name is already taken when uploading CSVs. (ebf4ab90c5db2589695eaf6c3d4c4206bad17e7b,
@epaulet)

## 1.9.0 - 14th September 2016
- Add `upload-csv` command for updating tabular variables without starting a new run.
(069e943cd94cbb08e6f00347ab6c8327372897ce, @epaulet)
- Add ability to filter uploads with `tag` and `site-id` and to upload specific test by file path.
(1d6e0a39d664e4c2d2135fa654571c28aaaac031, @epaulet)

## 1.8.1 - 19th August 2016
- Fix a bug that prevent uploading tests as a result from 1.8.0.
(651514ae94df6857c43e820fff60cdf8034f534d, @epaulet)

## 1.8.0 - 19th August 2016
- Filter RFML test downloads with `rainforest export` using `--tag`, `--folder`, and `--site-id` flags.
(5826000fddeb152dc1e2c8ad4baf04cdc0dd2001, @epaulet)

## 1.7.0 - 8th August 2016
- New run flag: `--wait` for hooking into existing run instead of starting a new one.
(77df41bf79b8635fb8c2d8a93968f975db092c69, @shosti)
- New run flag: `--junit-file` for exporting run results into JUnit format.
(349f2b1f5c8b423766875751c7cafed692fc2bed, @briancolby)
- New feature: Embed screenshots and file download links in your RFML tests.
(e081072c591c810b8bc3edce6d2e507d12a1a18e, @epaulet)

## 1.6.5 - 13th June 2016
- Exit with non-zero status code when `validate` fails. (8db1d38be39aa50d2afcdef817f78c654b3108b6,
@DyerSituation)
- Increase the amount of folders fetched from API to 100. (7fc7d426c029d35e73fe45712a79b35fda54ad60,
@DyerSituation)

## 1.6.4 - 10th June 2016
- Add redirects to non-embedded steps that need them as well. (2a5918d1a448f78016587c1711261e90a7be120f,
@epaulet)
- Fix a bug introduced in 1.6.1 where validating without an API token would result in an error.
(b604a5607e707057375ecc14db836d4b3ec537b1, @epaulet)

## 1.6.3 - 9th June 2016
- Add redirects to exported embedded tests. (ccbdaec2a9fb52c775ffbee05e203e1026b256f7, @theonegri)

## 1.6.2 - 8th June 2016
- Lower concurrent thread count to 16 and allow users to modify the amount of threads used.
(4333fca172e4e109a517ff6ffc5e11f89839fe85, @epaulet)

## 1.6.1 - 8th June 2016
- Use snappier API endpoint for validating RFML ids (e37a4e789eed9b329176ffdb16866aa8902a6cc5, @epaulet).

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
