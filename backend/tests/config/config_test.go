package config

import (
	"crawlab-lite/config"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInitConfig(t *testing.T) {
	Convey("Test InitConfig func", t, func() {
		x := config.InitConfig("../config.yml")

		Convey("The value should be nil", func() {
			So(x, ShouldEqual, nil)
		})
	})
}
