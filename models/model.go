package models

import (
	"gorm.io/gorm"
	"time"
)

type DANResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type LastResponse struct {
	CitizenLoginType int `json:"citizen_loginType,omitempty"`
	Services         struct {
		WS100101GetCitizenIDCardInfo struct {
			Request struct {
				Auth    interface{} `json:"auth"`
				CivilID interface{} `json:"civilId"`
				Regnum  string      `json:"regnum"`
			} `json:"request"`
			RequestID string `json:"requestId"`
			Response  struct {
				AddressApartmentName interface{} `json:"addressApartmentName"`
				AddressDetail        string      `json:"addressDetail"`
				AddressRegionName    interface{} `json:"addressRegionName"`
				AddressStreetName    string      `json:"addressStreetName"`
				AddressTownName      interface{} `json:"addressTownName"`
				AimagCityCode        string      `json:"aimagCityCode"`
				AimagCityName        string      `json:"aimagCityName"`
				BagKhorooCode        string      `json:"bagKhorooCode"`
				BagKhorooName        string      `json:"bagKhorooName"`
				BirthDate            string      `json:"birthDate"`
				BirthDateAsText      string      `json:"birthDateAsText"`
				BirthPlace           string      `json:"birthPlace"`
				CivilID              string      `json:"civilId"`
				Firstname            string      `json:"firstname"`
				Gender               string      `json:"gender"`
				Image                string      `json:"image"`
				Lastname             string      `json:"lastname"`
				Nationality          string      `json:"nationality"`
				PassportAddress      string      `json:"passportAddress"`
				PassportExpireDate   string      `json:"passportExpireDate"`
				PassportIssueDate    string      `json:"passportIssueDate"`
				PersonID             string      `json:"personId"`
				Regnum               string      `json:"regnum"`
				SoumDistrictCode     string      `json:"soumDistrictCode"`
				SoumDistrictName     string      `json:"soumDistrictName"`
				Surname              string      `json:"surname"`
			} `json:"response"`
			ResultCode    int    `json:"resultCode"`
			ResultMessage string `json:"resultMessage"`
		} `json:"WS100101_getCitizenIDCardInfo"`
	} `json:"services,omitempty"`
	Wsdl string `json:"wsdl,omitempty"`
}

type Users struct {
	Avatar         string         `gorm:"column:avatar" json:"avatar"`
	Bio            *string        `gorm:"column:bio" json:"bio"`
	Birthday       string         `gorm:"column:birthday" json:"birthday"`
	CreatedAt      *time.Time     `gorm:"column:created_at" json:"created_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at" json:"deleted_at"`
	Email          string         `gorm:"column:email" json:"email"`
	FcmToken       *string        `gorm:"column:fcm_token" json:"fcm_token"`
	FirstName      string         `gorm:"column:first_name" json:"first_name"`
	Gender         string         `gorm:"column:gender" json:"gender"`
	ID             int            `gorm:"column:id" json:"id"`
	LastName       string         `gorm:"column:last_name" json:"last_name"`
	Login          string         `gorm:"column:login" json:"login"`
	Password       string         `gorm:"column:password" json:"password"`
	Phone          string         `gorm:"column:phone" json:"phone"`
	RegisterNumber string         `gorm:"column:register_number" json:"register_number"`
	Hayag          string         `gorm:"column:hayag" json:"hayag"`
	Role           int            `gorm:"column:role" json:"role"`
	Status         *string        `gorm:"column:status" json:"status"`
	UpdatedAt      *time.Time     `gorm:"column:updated_at" json:"updated_at"`
}

func (u *Users) TableName() string {
	return "users"
}
