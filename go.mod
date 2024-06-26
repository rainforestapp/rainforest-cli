module github.com/rainforestapp/rainforest-cli

go 1.22

require (
	github.com/blang/semver v3.5.1+incompatible
	github.com/gyuho/goraph v0.0.0-20160328020532-d460590d53a9
	github.com/olekukonko/tablewriter v0.0.5
	github.com/rainforestapp/testutil v0.0.0-20170615220520-c9155e7da96e
	github.com/rhysd/go-github-selfupdate v1.2.3
	github.com/satori/go.uuid v1.2.0
	github.com/ukd1/go.detectci v0.0.0-20210512015713-5570f6cb6bb1
	github.com/urfave/cli v1.22.15
	github.com/whilp/git-urls v1.0.0
	golang.zx2c4.com/wireguard v0.0.0-20220920152132-bb719d3a6e2c
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20221104135756-97bc4ad4a1cb
	gvisor.dev/gvisor v0.0.0-20230927004350-cbd86285d259
	wiretap v0.0.0-00010101000000-000000000000
)

require (
	github.com/armon/go-socks5 v0.0.0-20160902184237-e75332964ef5 // indirect
	github.com/aws/aws-sdk-go v1.34.18 // indirect
	github.com/cpuguy83/go-md2man/v2 v2.0.4 // indirect
	github.com/garyburd/redigo v1.6.2 // indirect
	github.com/go-ping/ping v1.1.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/go-github/v30 v30.1.0 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/google/gopacket v1.1.19 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/inconshreveable/go-update v0.0.0-20160112193335-8152e7eb6ccf // indirect
	github.com/jmespath/go-jmespath v0.3.0 // indirect
	github.com/libp2p/go-reuseport v0.2.0 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/tcnksm/go-gitconfig v0.1.2 // indirect
	github.com/ulikunitz/xz v0.5.9 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.17.0 // indirect
	golang.org/x/oauth2 v0.4.0 // indirect
	golang.org/x/sync v0.3.0 // indirect
	golang.org/x/sys v0.15.0 // indirect
	golang.org/x/time v0.1.0 // indirect
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.2-0.20230118093459-a9481185b34d // indirect
	gopkg.in/check.v1 v1.0.0-20200902074654-038fdea0a05b // indirect
)

replace github.com/rhysd/go-github-selfupdate => github.com/rainforestapp/go-github-selfupdate v1.2.4-0.20210729013827-905f4fc54255

// For Rainforest Direct Connect, we depend on wiretap
replace wiretap => ./dependencies/wiretap/src

// Wiretap requires a fork of wireguard-go, and even though wiretap specifies this replace directive in it's go.mod, we have to duplicate it here
replace golang.zx2c4.com/wireguard => github.com/luker983/wireguard-go v0.0.0-20231019223227-fc689040dc0a
