package routes

import (
	"bytes"
	"crawlab-lite/config"
	"crawlab-lite/controllers"
	"crawlab-lite/lib/validate_bridge"
	"crawlab-lite/routes"
	"encoding/json"
	"github.com/apex/log"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"net/http"
)

func InitTestApp() *gin.Engine {
	binding.Validator = new(validate_bridge.DefaultValidator)
	app := gin.Default()
	if err := config.InitConfig("../../config.yml"); err != nil {
		log.Error("Init config error:" + err.Error())
		panic(err)
	}
	routes.InitRoutes(app)
	return app
}

func PostJson(url string, jsonBody map[string]string) *http.Request {
	req, _ := http.NewRequest("POST", url, JsonBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}

func JsonBuffer(m map[string]string) *bytes.Buffer {
	jsonValue, _ := json.Marshal(m)
	return bytes.NewBuffer(jsonValue)
}

func GetResponse(body *bytes.Buffer) *controllers.Response {
	var resp *controllers.Response
	_ = json.Unmarshal(body.Bytes(), &resp)
	return resp
}
