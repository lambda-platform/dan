_# DAN e-mongolia for Lambda Platform
Config Examples

1. Set .env client information

DAN_REDIRECT_URL="https://xxxx.mn/sso"
DAN_REDIRECT_ROUTE="/sso"
DAN_CLIENT_ID=XXXXX
DAN_CONSUMER_SECRET=XXXXX

2.import dan in bootstrap/bootstrap.go

`import "github.com/lambda-platform/dan"`

3. set dat in in bootstrap/bootstrap.go - >func Set() *lambda.Lambda {

`func Set() *lambda.Lambda {
.
.
.
dan.Set(Lambda.App)
}
`
