Project listing the services (old-named "module") in your current appengine project

# Go

Use the `go@1.9` appengine version.

# Install dependencies

```
go get golang.org/x/oauth2
go get google.golang.org/appengine
go get cloud.google.com/go
```

# Test

`goapp serve .`

# Deploy

Not working with the following command:
`gcloud app deploy --application XX ./app.yaml --version go-listing`
