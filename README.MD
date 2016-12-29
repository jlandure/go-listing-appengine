Project listing the services (old-named "module") in your current appengine project

# Test

`goapp serve .`

# Deploy

`goapp deploy -application my-project -version go-listing .`

# Problem

Not working with the following command:
`gcloud app deploy ./app.yaml --version go-listing --quiet`