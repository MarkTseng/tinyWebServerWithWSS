package main

import (
	"encoding/json"
	"fmt"
	limit "github.com/bu/gin-access-limit"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/contrib/secure"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const (
	configFile        = "config.json"
	WWW_ROOTDIR_ROUTE = "/"
)

// save server startup parameter
var ConfigSetting map[string]interface{}

func configParse() {
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	if err := json.Unmarshal(content, &ConfigSetting); err != nil {
		panic(err)
	}
}

func ensureLoggedIn() gin.HandlerFunc {
	return func(c *gin.Context) {
		loggedInInterface, _ := c.Get("is_logged_in")
		if loggedInInterface != nil {
			loggedIn := loggedInInterface.(bool)
			if !loggedIn {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		}
	}
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if token, err := c.Cookie("token"); err == nil || token != "" {
			c.Set("is_logged_in", true)
		} else {
			c.Set("is_logged_in", false)
			if !strings.Contains(c.Request.URL.Path, "login") {
				c.AbortWithStatus(http.StatusUnauthorized)
			}
		}
	}
}
func generateSessionToken() string {
	return uuid.New()
}

func tokenkeyCheck(c *gin.Context) {
	tokenkey := c.Query("tokenkey") // shortcut for c.Request.URL.Query().Get("lastname")
	if tokenkey == ConfigSetting["TOKENKEY"] {
		// If the username/password is valid set the token in a cookie
		sessionToken := generateSessionToken()
		cookieExpireTime, _ := strconv.Atoi(ConfigSetting["COOKIE_EXPIRE_TIME"].(string))
		c.SetCookie("token", sessionToken, cookieExpireTime, "", "", false, true)
		c.Set("is_logged_in", true)

		session := sessions.Default(c)
		session.Set("login", "true")
		err := session.Save()
		if err != nil {
			log.Println("user session svae failed")
		}
		c.Redirect(http.StatusTemporaryRedirect, WWW_ROOTDIR_ROUTE)
	} else {
		c.String(http.StatusUnauthorized, "tokenkey %s is not correct", tokenkey)
	}
}

func main() {
	fmt.Println("tiny Web Server by Mark Tseng")
	configParse()
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	//r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// set source connect IP/netmask
	if ConfigSetting["CLIENT_NETWORK_LIMIT"] != nil {
		r.Use(limit.CIDR(ConfigSetting["CLIENT_NETWORK_LIMIT"].(string)))
	}
	// Set sessions for keeping user info
	store := sessions.NewCookieStore([]byte("C4XsecretSession"))
	r.Use(sessions.Sessions("C4XSession", store))

	// set secure options and static html
	r.Use(AuthRequired(),
		secure.Secure(secure.Options{
			SSLRedirect:          true,
			SSLProxyHeaders:      map[string]string{"X-Forwarded-Proto": "https"},
			STSSeconds:           315360000,
			STSIncludeSubdomains: true,
			FrameDeny:            true,
			ContentTypeNosniff:   true,
			BrowserXssFilter:     true,
		}),
		static.Serve(WWW_ROOTDIR_ROUTE, static.LocalFile("public", true)))
	r.GET("/login", tokenkeyCheck) //login?tokenkey=1234567890
	authorized := r.Group(WWW_ROOTDIR_ROUTE)
	authorized.Use(AuthRequired())

	r.RunTLS(fmt.Sprintf("%s:%s", ConfigSetting["SERVER_IP"], ConfigSetting["SERVER_PORT"]), "ssldata/certificate.crt", "ssldata/private.key")
}
