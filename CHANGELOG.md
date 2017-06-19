# Rainforest CLI Changelog

## 2.1.0 - 19 Jun 2017
- Added run group support for future run group feature.
(86a4573db19cb2b5aef7a53c765d0121be60520f, @sondhayni-rfqa)
(21e9fda469a23f40a9b208e8660b4b2b80d00c86, @epaulet)
- Replace all non-alphanumeric characters with underscores when creating RFML
files. (21e9fda469a23f40a9b208e8660b4b2b80d00c86, @epaulet)

## 2.0.4 - 5 Jun 2017
- Log errors when attempting to upload tests with embedded files that do not
exist locally, but upload the test anyway. This behavior is backwards compatible
with versions 1.X.
(52cf356f6d4a1d4359537e53923949facd5d5c08, @epaulet)

## 2.0.3 - 2 Jun 2017
- You may now either omit the browsers attribute or leave the browser list
empty to set the default browsers for a test as none.
(49d48abf5b6c3591f6998622a34884426d9526a1, @epaulet)

## 2.0.2 - 2nd May 2017
- Replace illegal file path characters when creating RFML files.
(175c98e6568a909cd9a000a8381768d7189aa25a, @epaulet)

## 2.0.1 - 12th April 2017
- Download all tests from test API and return proper errors.
(726f2de5215d66eeb76aa530f76b4a8a59e76f71, @epaulet)
