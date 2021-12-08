module github.com/rainforestapp/rainforest-cli

go 1.12

require (
	github.com/aws/aws-sdk-go v1.34.18 // indirect
	github.com/blang/semver v3.5.1+incompatible
	github.com/garyburd/redigo v1.6.2 // indirect
	github.com/gyuho/goraph v0.0.0-20160328020532-d460590d53a9
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/olekukonko/tablewriter v0.0.0-20160621093029-daf2955e742c
	github.com/rainforestapp/testutil v0.0.0-20170615220520-c9155e7da96e
	github.com/rhysd/go-github-selfupdate v1.2.3
	github.com/satori/go.uuid v1.2.0
	github.com/ukd1/go.detectci v0.0.0-20210512015713-5570f6cb6bb1
	github.com/urfave/cli v1.22.5
	github.com/whilp/git-urls v1.0.0
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
)

replace github.com/rhysd/go-github-selfupdate => github.com/rainforestapp/go-github-selfupdate v1.2.4-0.20210729013827-905f4fc54255
