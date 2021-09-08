# Rainforest CLI Changelog

## 2.21.4 - 2021-09-07
- Fix for passed / failed reporting being in the wrong order 🤦‍♂️
  - (504c1c129fe1847d85fc6e07f922352f8d9679e4, @ukd1)

## 2.21.3 - 2021-09-07
- Expand the output of the CLI a little, making it easier to get to Rainforest and see progress.
  - (9f932b75eecc33e3bada7e854022913bcedeed9a, @ukd1)

## 2.21.2 - 2021-08-30
- Support for cleaning up of temporary environments
  - (ba1846611d08f824f5bd6f4174d9f0ccac725409, @ukd1)

## 2.21.1 - 2021-08-12
- Stop publishing releases to Equinox.io
  - (f26fb192f4f4ed957e069582f8eb49b916202fb1, @magni-)

## 2.21.0 - 2021-08-09
- Switch autoupdate functionality provider from Equinox.io to `go-github-selfupdate`
  - (a5937b1ac787063dcdafea483967a81e551e2882, @magni-)

## 2.20.0 - 2021-07-29
- Include URLs to test results for failures
  - (5d09380, @jbarber)

## 2.19.1 - 2021-07-14
- Publish releases to GitHub via GoReleaser
  - (ab5f381e0d7c4e2594c37718846fe90d3a76d4c0, @magni-)

## 2.19.0 - 2021-05-14
- Update docs, add telemetry + customization based on your CI type & git settings
  - (56a02b6, @ukd1)

## 2.18.2 - 2021-03-16
- Update --custom-url documentation
  - (e492917, @jbarber)

## 2.18.1 - 2021-01-27
- Find the CLI's path dynamically for re-exec'ing when using the --max-reruns flag
  - (4f79806, @jbarber)

## 2.18.0 - 2021-01-15
- Add `--max-reruns N` flag, which allows re-trying failed tests in a run up to `N` times.
  - (b22c33191c5e36342840f5aa22f802083f8eaca9, @maciejgryka)

## 2.17.0 - 2020-09-08
- Preserve draft state when downloading RFMLs
  - (f904bdb68de60e2ff213c3d37f7ab7af20c19745, @mikesbrown)

## 2.16.1 - 2020-07-14
- Fix default crowd setting being overridden for Run Group runs
  - (7d88de2155220694fc1423f9f32df8ee6a2dc59b, @magni-)
- Fix publishing to secondary Equinox channels when releasing a new stable or beta version
  - (c8c38d75ee6a1ec37098583e70253ecd170b0ad0, @magni-)
- Fix GCR image triggering
  - (05da067539ddce74f2f00729f7ccf7146b77795d, @AntGH)

## 2.16.0 - 2020-07-08
- Add [`rerun`](README.md#running-tests) command
  - (b54994797897e98a724a5ae330f9bf0faa88417c, @magni-)

## 2.15.3 - 2020-05-21
- Add rainforest-orb and version to UserAgent if CLI is called from CircleCI Orb
  - (da19bb6ad2702aead89f7169504af7286659d621, @magni-)

## 2.15.2 - 2020-03-20
- Report Rainforest Automation failures in JUnit
  - (2ceb9f41fac2cf7136376e0e9ba3f79ff011b611, @kruszczynski)

## 2.15.1 - 2019-11-26
- Remove `skip_mark_as_viewed` param when fetching run results
  - (8985398c8e0244241421bdf2a56a57cc93c8e1db, @AJFunk)

## 2.15.0 - 2019-11-18
- Support `automation_and_crowd` option in the `--crowd` flag.
  - (5835fb69e68e8cd57cd259bc3f68400f63be1789, @kruszczynski)

## 2.14.0 - 2019-07-19
- Support `automation` option in the `--crowd` flag.
  - (febca56232f57397c1b95df4478b3e2e58a2a415, @grodowski)

## 2.13.0 - 2019-05-24
- Add `priority` setting to RFML.
  - (db095be3f3463e1af7aa829f48ecfa258df38c44, @fuzzylita)

## 2.12.0 - 2019-04-25
- Add a `release` flag to the `run` command.
  - (cbb6289bdd54cb3f294819621e9a8bcdc7d9cdd0, @AJFunk, @fuzzylita)

## 2.11.3 - 2019-04-23
- Increase upload app slots to 100 for `mobile-upload` command `--app-slot` flag.
  - (abc9966adb6364d567b75b9c11ae53fbabd6edc3, @shanempope)

## 2.11.2 - 2019-04-12
- Add multiple app support to `mobile-upload` command via `--app-slot` flag.
  - (2374999d7965bb5854a15b6ad1d885922ec788f5, @shanempope)

## 2.11.1 - 2019-03-19
- Fix parsing errors for non-JSON responses.
  - (ec4ee32934c670a9a1b285f0eef3775cfab91de5, @epaulet)

## 2.11.0 - 2019-03-15
- Add `mobile-upload` command.
  - (ccb84a2a923c7e4a142a1d0a3f4d475bb66127e4, @shanempope)
- Add `environments` command.
  - (f6f53c0971aeccad4e6674bd9965da72b4a8fa43, @shanempope)

## 2.10.3 - 2019-01-16
- Remove local validation for RFML state attribute so that the API will validate it instead.
  - (ac67e3d4c3f6f7d8df2ea465c3741980dfc16af2, @epaulet, @jbarber)

## 2.10.2 - 2019-01-14
- Use `include_feedback` and `skip_mark_as_viewed` query params when fetching run results in order to get details of failures and not update the results state to `viewed`.
  - (c76cb1dea35677f951dc22bc1b3ed3b7895bb569, @epaulet)
- Remove extra request to API that was used for fetching the `updated_at` attribute on test results.
  - (053157a38cd35be14a0fe371711e03a03a179295, @epaulet)

## 2.10.1 - 2018-10-04
- Use slim=true when interacting with the tests API endpoint.
  - (54e7a8918c95b3b71d5ef91d007fa957e4b3a32b, @nxvl)

## 2.10.0 - 2018-08-15
- Remove defunct `--embed-tests` flag and replace with `--flatten-steps` flag.
  - (804b18864407cf3cad6b04fec0a4f8d5d67f0a00, @epaulet)

## 2.9.0 - 2018-05-24
- Show more descriptive Rainforest API errors.
  - (071123b4e9301357b7df6dda78dd79229c953530, @epaulet, @jbarber)
  - (778fbac8a8130f9422baade38eda25f49f879712, @epaulet)
- Use `--debug` flag to print raw response bodies from Rainforest API in case
  the response body is not parseable JSON.
  - (071123b4e9301357b7df6dda78dd79229c953530, @epaulet)

## 2.8.12 - 2018-04-19
- Fix bug with `rainforest rm` failing silently if the file cannot be parsed.
  - (aca9793b49de56bb4b2cc1770950daffaa709391, @epaulet)
- Fix grammatical errors.
  - (64f578845f498ef06beca3edacbc9f4a3fd2f09f, @jeis2497052)a

## 2.8.11 - 2017-12-07
- Patch over bug with parsing `null` timestamps returned by Rainforest generator API.
  - (c8be6923d987522b08ea9b20088a3ab7f14e1a50, @epaulet)

## 2.8.10 - 2017-11-29
- Fix bug with `-f` flag returning an odd error.
  - (c9242ba02b13e439203499cc3e965ed3c89c82d4, @macocha)

## 2.8.9 - 2017-11-06
- Re-release 2.8.7 with fixes
  - (58e3464f8f7806db2830ed9411a27f87d36f1dbe, @epaulet)

## 2.8.8 - 2017-11-03
- Revert 2.8.7
  - (a1d29d8b253259aee1eb6906510492a6f461f63e, @jbarber)

## 2.8.7 - 2017-11-02
- Properly exit with an error status when using flags improperly.
  - (8ae777dce7d1b6f9abdf85bf353aecfc6e38be17, @epaulet)

## 2.8.6 - 2017-11-01
- Remove duplicate listing of `run-groups` option.
  - (0d54c1eb5d4dff75eb3b38ce67c7594c6788918c, @epaulet)

## 2.8.5 - 2017-10-30
- Fix error that occurs when `-f` if given as the final option of a command.
  - (d0abed0156ccc1da184a73fe4c0f6b2d5b6cce91, @epaulet)

## 2.8.4 - 2017-10-26
- Default the value of the local file flag to false when not given.
  - (1222cc0acb99995d0c29aa8d1acee400b4923af1, @epaulet)

## 2.8.3 - 2017-10-25
- Correctly find and parse failed steps for JUnit reports.
  - (090e407a579e452d1c421a36c531cdd4de4bd60d, @epaulet)

## 2.8.2 - 2017-10-13
- Add a default JUnit test suite name for runs without a description.
  - (b66a7aaeb952986f29303f0291e4c749d990bd36, @epaulet)

## 2.8.1 - 2017-10-10
- Add file path to invalid file error message when using the `validate` command.
  - (be557fff7e39c478426f3ecd6289fe5689c7a83f, @epaulet)

## 2.8.0 - 2017-10-03
- Add `feature_id` and `state` fields to RFML headers.
  - (7e6b0b87e438838cef7b28fda4396f739f819c33, @epaulet)

## 2.7.2 - 2017-09-28
- Fix issue with process not exiting when tracking a complete run
  - (135a2f0ccd96ca3f9ed32bd6a42251021fc11603, @epaulet)

## 2.7.1 - 2017-09-28
- Rerelease in attempt to revert to 2.5.0 again.

## 2.7.0 - 2017-09-28
- Revert previous release as it never stopped trying
  - (fd149aebe61e89395c7f9673591f00c672f5e4f1, @jbarber)

## 2.6.0 - 2017-09-28
- When using --wait don't give up if the API returns an error
  - (95eb3d38dcf9f8cf4232c34f15ef5086bb51c9a9, @jbarber)

## 2.5.0 - 2017-09-12
- Improve run group support: runs started from run groups will now apply run group browser settings.
- Add support for viewing and filtering by features
  - (e4fe58df872c178fd39756424983da33e4dd96a0, @shosti)
- Fix a bug in printing run groups
  - (32e0cd770e9a3faccce6828f3970d4a83181af6b, @shosti)

## 2.4.0 - 2017-09-03
- Exit with non-zero status if an unknown flag is given
  - (2ebdd906d8e4314b3c0db4d0a72a2d6ca2af52ee, @jbarber)

## 2.3.0 - 2017-09-01
- Fix bug that parsed a remote file reference as a file path when using file.screenshot and file.download step variables.
  - (6a443d7e0ef5cc1ce9696ee07552ed87611e450b, @epaulet)
- Add new prerelease feature for running local RFML tests.
  - (573fdb7f179aa12baf99a4b2bf351649633d1636, @shosti)

## 2.2.0 - 2017-07-14
- Added a `--debug` flag to print out HTTP headers.
  - (a9bc9dde31124f1f37934c8e85c4bd11692a8f9c, @sondhayni-rfqa)

## 2.1.1 - 2017-06-22
- Changed default file name when downloading RFML tests:
  - Do not use sequences of multiple underscores in file name.
    - (12f2dafa9cd62d489c4837055535fb45580b9ef8, @epaulet)
  - Do not use more than 30 characters from a test's title in file name.
    - (c2f1baeafff3f89eff3f32295f301f9be6211dda, @epaulet)

## 2.1.0 - 2017-06-19
- Added run group support for future run group feature.
  - (86a4573db19cb2b5aef7a53c765d0121be60520f, @sondhayni-rfqa)
  - (21e9fda469a23f40a9b208e8660b4b2b80d00c86, @epaulet)
- Replace all non-alphanumeric characters with underscores when creating RFML files.
  - (21e9fda469a23f40a9b208e8660b4b2b80d00c86, @epaulet)

## 2.0.4 - 2017--06-05
- Log errors when attempting to upload tests with embedded files that do not exist locally, but upload the test anyway. This behavior is backwards compatible with versions 1.X.
  - (52cf356f6d4a1d4359537e53923949facd5d5c08, @epaulet)

## 2.0.3 - 2017-06-02
- You may now either omit the browsers attribute or leave the browser list empty to set the default browsers for a test as none.
  - (49d48abf5b6c3591f6998622a34884426d9526a1, @epaulet)

## 2.0.2 - 2017-05-02
- Replace illegal file path characters when creating RFML files.
  - (175c98e6568a909cd9a000a8381768d7189aa25a, @epaulet)

## 2.0.1 - 2017-04-12
- Download all tests from test API and return proper errors.
  - (726f2de5215d66eeb76aa530f76b4a8a59e76f71, @epaulet)
