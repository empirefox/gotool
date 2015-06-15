package ng

import (
	"bytes"
	"io/ioutil"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type instance struct {
	Name string `json:"name"`
	Age  int    `json:"age,omitempty"`
}

func TestWrite(t *testing.T) {
	Convey("Write", t, func() {
		Convey("should write ng module to Writer", func() {
			buf := &bytes.Buffer{}
			m := Module{
				Type:       "constant",
				ModuleName: "app.constants.test",
				Name:       "Tester",
				Instance:   instance{Name: "myname"},
			}
			err := Write(buf, m)
			So(err, ShouldBeNil)

			raw, err := ioutil.ReadAll(buf)
			So(err, ShouldBeNil)

			So(string(raw), ShouldEqual, `angular.module('app.constants.test', []).constant('Tester', {"name":"myname"});`)
		})
	})
}
