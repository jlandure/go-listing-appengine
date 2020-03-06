Project listing the services in your current appengine project

# Go

Use the `go 1.13` version.

# Install gcloud command

Use the `gcloud components install app-engine-go` command

# Install dependencies

```
go mod vendor
```

# Documentation

If needed, please read these links:

- [Admin API](https://cloud.google.com/appengine/docs/admin-api/reference/rest/v1/apps.services/list)
- [Go Doc for google.golang.org/api/appengine/v1](https://pkg.go.dev/google.golang.org/api/appengine/v1?tab=doc)

# Test

`go run listing.go`

# Deploy

Not working with the following command:
`gcloud app deploy --project XX ./app.yaml --version go-listing`
