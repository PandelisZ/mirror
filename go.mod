module github.com/pandelisz/mirror

go 1.14

require (
	github.com/google/go-github/v32 v32.0.0
	github.com/kr/pretty v0.1.0 // indirect
	github.com/libgit2/git2go/v30 v30.0.5
	github.com/tj/go-spin v1.1.0
	github.com/urfave/cli/v2 v2.2.0
	github.com/xanzy/go-gitlab v0.33.0
	golang.org/x/crypto v0.0.0-20190530122614-20be4c3c3ed5 // indirect
	golang.org/x/oauth2 v0.0.0-20181106182150-f42d05182288
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
)

// replace github.com/libgit2/git2go/v30 ../../libgit2/git2go
