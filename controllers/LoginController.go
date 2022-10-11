package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lambda-platform/dan/models"

	"github.com/lambda-platform/lambda/DB"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"time"
)

var clientID string = os.Getenv("DAN_CLIENT_ID")
var clientSecret string = os.Getenv("DAN_CONSUMER_SECRET")
var redirectURI string = os.Getenv("DAN_REDIRECT_URL")

func DANRedirect(c *fiber.Ctx) error {

	state := GenerateSecureToken(12)
	scope := "WwogIHsKICAgICJzZXJ2aWNlcyI6IFsKICAgICAgIldTMTAwMTAxX2dldENpdGl6ZW5JRENhcmRJbmZvIgogICAgXSwKICAgICJ3c2RsIjogImh0dHBzOi8veHlwLmdvdi5tbi9jaXRpemVuLTEuMy4wL3dzP1dTREwiCiAgfQpd"
	DAN := fmt.Sprintf("https://sso.gov.mn/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s&scope=%s&state=%s", clientID, redirectURI, scope, state)

	fmt.Println(DAN)
	fmt.Println(scope)

	return c.Status(http.StatusSeeOther).Redirect(DAN)

}

func AuthWithDan(c *fiber.Ctx) error {
	code := c.Query("code")

	if code != "" {

		response := models.DANResponse{}

		url := "https://sso.gov.mn/oauth2/token"
		method := "POST"

		payload := &bytes.Buffer{}
		writer := multipart.NewWriter(payload)
		_ = writer.WriteField("grant_type", "authorization_code")
		_ = writer.WriteField("code", code)
		_ = writer.WriteField("client_id", clientID)
		_ = writer.WriteField("client_secret", clientSecret)
		_ = writer.WriteField("redirect_uri", redirectURI)
		//_ = writer.WriteField("scope", scope)
		//_ = writer.WriteField("login_type", "OTP")
		err := writer.Close()
		if err != nil {

			return err
		}

		client := &http.Client{}
		req, err := http.NewRequest(method, url, payload)

		if err != nil {

			return err
		}

		req.Header.Set("Content-Type", writer.FormDataContentType())
		res, err := client.Do(req)
		if err != nil {

			return err
		}
		defer res.Body.Close()

		body, err := ioutil.ReadAll(res.Body)
		if err != nil {

			return err
		}
		fmt.Println(string(body))
		json.Unmarshal(body, &response)

		if response.AccessToken != "" {

			danresponse, errDan := GetCitizenData(response.AccessToken)

			if errDan != nil {
				return errDan
			}

			if len(danresponse) >= 2 {
				fmt.Println(danresponse[1].Services.WS100101GetCitizenIDCardInfo.ResultCode)
				fmt.Println(danresponse[1].Services.WS100101GetCitizenIDCardInfo.ResultMessage)
				if danresponse[1].Services.WS100101GetCitizenIDCardInfo.ResultCode == 0 && danresponse[1].Services.WS100101GetCitizenIDCardInfo.ResultMessage == "амжилттай" {
					return DANSUCESS(c, danresponse)
				}
			} else {
				fmt.Println(danresponse)
			}

			return c.Status(http.StatusFound).Redirect("/")
		} else {
			return c.Status(http.StatusFound).Redirect("/")
		}

	} else {
		return c.Status(http.StatusFound).Redirect("/")
	}

}
func DANSUCESS(c *fiber.Ctx, danresponse []models.LastResponse) error {

	user := danresponse[1].Services.WS100101GetCitizenIDCardInfo.Response
	fmt.Println(user.CivilID)
	fmt.Println(user.Regnum)
	foundUser := agentUtils.AuthUserObjectByLogin(user.CivilID)

	var roleID int64 = 0
	var userID int64 = 0

	if len(foundUser) == 0 {
		newUser := models.Users{}
		newUser.Role = 8
		newUser.FirstName = user.Firstname
		newUser.LastName = user.Lastname
		newUser.Gender = "m"
		if user.Gender != "Эрэгтэй" {
			user.Gender = "f"
		}
		newUser.Login = user.CivilID
		newUser.RegisterNumber = user.Regnum
		newUser.Email = user.Regnum
		newUser.Hayag = user.PassportAddress
		newUser.Birthday = user.BirthDateAsText[0:10]

		dec, err := base64.StdEncoding.DecodeString(user.Image)
		if err != nil {
			panic(err)
		}

		if _, err := os.Stat("public/uploaded/dan/"); os.IsNotExist(err) {
			os.MkdirAll("public/uploaded/dan/", 0755)
		}
		imageName := "/uploaded/dan/" + user.CivilID + ".png"
		f, err := os.Create("public" + imageName)
		if err != nil {
			panic(err)
		}
		defer f.Close()

		if _, err := f.Write(dec); err != nil {
			panic(err)
		}
		if err := f.Sync(); err != nil {
			panic(err)
		}
		newUser.Avatar = imageName

		DB.DB.Create(&newUser)

		foundUser = agentUtils.AuthUserObjectByLogin(user.CivilID)
	}

	if reflect.TypeOf(foundUser["id"]).String() == "string" {
		i, err := strconv.ParseInt(foundUser["id"].(string), 10, 64)
		if err != nil {
			panic(err)
		}
		userID = i

	} else {
		userID = foundUser["id"].(int64)
	}

	token, _ := createJwtToken(UserData{Id: userID, Login: foundUser["login"].(string), Role: roleID})

	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Path = "/"
	cookie.Value = token
	cookie.Expires = time.Now().Add(time.Hour * time.Duration(config.Config.JWT.Ttl))

	delete(foundUser, "password")

	foundUser["jwt"] = token
	userString, _ := json.Marshal(foundUser)
	c.Set("Content-Type", "text/html")
	return c.SendString(fmt.Sprintf(`<!doctype html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport"
          content="width=device-width, user-scalable=no, initial-scale=1.0, maximum-scale=1.0, minimum-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>DAN</title>
</head>
<body>
<script>

    localStorage.setItem('user', '%s');
    localStorage.setItem('token', '%s');
    localStorage.setItem('user_token', '%s');

    window.location.replace("/");
</script>
</body>
</html>
`, string(userString), token, token))
}
func GetCitizenData(accessToken string) ([]models.LastResponse, error) {
	url := "https://sso.gov.mn/oauth2/api/v1/service"
	method := "GET"
	LastResponse := []models.LastResponse{}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return LastResponse, err
	}
	req.Header.Add("authorization", "Bearer "+accessToken)
	req.Header.Add("content-type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return LastResponse, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {

		return LastResponse, err
	}

	fmt.Println(string(body))

	errGet := json.Unmarshal(body, &LastResponse)

	if errGet != nil {
		return LastResponse, errGet
	}

	return LastResponse, nil

}
func GenerateSecureToken(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
func createJwtToken(user UserData) (string, error) {
	// Set custom claims
	claims := &jwtClaims{
		user.Id,
		user.Login,
		user.Role,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(config.Config.JWT.Ttl)).Unix(),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	t, err := token.SignedString([]byte(config.Config.JWT.Secret))
	if err != nil {
		return "", err
	}
	return t, nil
}

type UserData struct {
	Id    int64
	Login string
	Role  int64
}
type jwtClaims struct {
	Id    int64  `json:"id"`
	Login string `json:"login"`
	Role  int64  `json:"role"`
	jwt.StandardClaims
}
