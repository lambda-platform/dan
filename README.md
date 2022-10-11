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

3. Create user model in app/models/user.go
`package models

import (
"github.com/lambda-platform/lambda/DB"
"gorm.io/gorm"
"time"
)

type Users struct {

Avatar                *string        `gorm:"column:avatar" json:"avatar"`
Bio                   *string        `gorm:"column:bio" json:"bio"`
Birthday              *DB.Date       `gorm:"column:birthday" json:"birthday"`
CreatedAt             *time.Time     `gorm:"column:created_at" json:"created_at"`
DeletedAt             gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
Email                 string         `gorm:"column:email" json:"email"`
FcmToken              *string        `gorm:"column:fcm_token" json:"fcm_token"`
FirstName             string         `gorm:"column:first_name" json:"first_name"`
Gender                string         `gorm:"column:gender" json:"gender"`
ID                    int            `gorm:"column:id" json:"id"`
LastName              string         `gorm:"column:last_name" json:"last_name"`
Login                 string         `gorm:"column:login" json:"login"`
Password              string         `gorm:"column:password" json:"password"`
Phone                 string         `gorm:"column:phone" json:"phone"`
RegisterNumber        string         `gorm:"column:register_number" json:"register_number"`
Role                  *int           `gorm:"column:role" json:"role"`
Status                *string        `gorm:"column:status" json:"status"`
UpdatedAt             *time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

func (u *Users) TableName() string {
return "users"
}

`
