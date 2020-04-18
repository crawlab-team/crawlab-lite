package routes

import (
	"crawlab-lite/model"
	. "github.com/smartystreets/goconvey/convey"
	"net/http/httptest"
	"testing"
)

func TestUserRoutes(t *testing.T) {
	Convey("Test User Routes", t, func() {
		app := InitTestApp()
		users, err := model.GetUserList()
		So(err, ShouldBeNil)
		user := users[0]

		Convey("Test correct user", func() {
			w := httptest.NewRecorder()
			values := map[string]string{"username": user.Username, "password": user.Password}
			req := PostJson("/api/login", values)
			app.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, 200)

			resp := GetResponse(w.Body)
			So(resp.Code, ShouldEqual, 200)
			So(resp.Message, ShouldEqual, "success")
			So(resp.Data, ShouldNotBeEmpty)
		})

		Convey("Test wrong user", func() {
			w := httptest.NewRecorder()
			values := map[string]string{"username": "abcdefg", "password": "000000"}
			req := PostJson("/api/login", values)
			app.ServeHTTP(w, req)

			So(w.Code, ShouldEqual, 401)

			resp := GetResponse(w.Body)
			So(resp.Code, ShouldEqual, 401)
			So(resp.Message, ShouldEqual, "not authorized")
			So(resp.Data, ShouldBeEmpty)
		})
	})
}
