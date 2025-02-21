package psutil_test

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/shirou/gopsutil/v4/sensors"
)

func TestPsutil(t *testing.T) {
	is := is.New(t)
	sensors, err := sensors.SensorsTemperatures()
	is.NoErr(err)
	for _, s := range sensors {
		fmt.Printf("%#v\n\n", s)
	}
}
