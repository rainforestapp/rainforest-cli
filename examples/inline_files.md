## Embedding Files Inline

You can embed inline screenshots and file download links directly into RFML by simply supplying the correct step variable and file path in your step.

The file paths used by `file.screenshot` and `file.download` are _relative to the RFML file_. **Your
file must also be checked into version control.** This ensures that everyone always has the most
up to date file when using the `upload` command.

#### Inline Screenshots Example

In the example below if your RFML test is located at `~/Desktop/my_test_repo/my_awesome_test.rfml`, then the screenshot in question would be located at `~/Desktop/my_test_repo/test_screenshot.png`.

```
#! my_rfml_id
# title: My Awesome Test
# start_uri: /
#

Click on the button in this image {{ file.screenshot(./test_screenshot.png) }}.
Did you go to a login page?
```

#### Inline Downloadable File Example
In the example below if your RFML test is located at `~/Desktop/my_test_repo/my_awesome_test.rfml`, then the file in question would be located at `~/Desktop/my_test_repo/test_file.txt`.

```
#! my_rfml_id
# title: My Awesome Test
# start_uri: /
#

Click on the button in this image {{ file.download(./test_file.txt) }}.
Did you go to a login page?
```


#### Uploading

Once you have set up your variables, you're done! `rainforest upload` will take care of the rest. You can double check that your uploads were successful but making sure that the file paths were replaced with proper arguments in the dashboard. Proper arguments look similar to the following:

**Screenshots:** `{{ file.screenshot(2262, WVXkKK) }}`

**Downloads:** `{{ file.download(2262, WVXkKK, test_screenshot.png) }}`
