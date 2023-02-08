module NBServer/tools/deploymentgen

go 1.18

// 替换规则
replace github.com/shafreeck/cortana => github.com/liwenson/cortana v0.0.0-20221028085157-d44dd63d766a

require (
	github.com/flosch/pongo2/v6 v6.0.0
	github.com/ghodss/yaml v1.0.0
	github.com/shafreeck/cortana v0.0.0-20220426120812-6e06d65fed3b
	github.com/xanzy/go-gitlab v0.73.1
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.1 // indirect
	golang.org/x/net v0.0.0-20220805013720-a33c5aa5df48 // indirect
	golang.org/x/oauth2 v0.0.0-20220722155238-128564f6959c // indirect
	golang.org/x/time v0.0.0-20220722155302-e5dcc9cfc0b9 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)
