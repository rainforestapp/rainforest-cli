## List of Notable changes from 1.X to 2.0:

- Global flags are given before the command. eg:
    - 1.X: `rainforest run --token my_token --tag foo`
    - 2.0: `rainforest --token my_token run --tag foo`
- `--fg` has been removed. Runs are tracked in the foreground by default. To create
a run without tracking it, use the `--bg` flag.
- The `export` command has been renamed `download`.
- `--import-variable-csv-file` has been deprecated for the `csv-upload` command.
Instead, simply supply the file path as the argument without a flag.
    - Note: `--import-variable-csv-file` with the `run` command has not been
    deprecated.
- New global flag `--skip-update` will execute your command without attempting
to auto-update your CLI binary if an update is available.
