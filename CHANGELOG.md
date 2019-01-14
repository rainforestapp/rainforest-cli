# Rainforest CLI Changelog

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
