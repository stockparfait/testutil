// Copyright 2022 Stock Parfait

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

//     http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testutil

import (
	"math"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestRound(t *testing.T) {
	t.Parallel()

	Convey("Rounding functions work", t, func() {
		Convey("RoundFixed", func() {
			So(RoundFixed(0.017, 2), ShouldEqual, 0.02)
			So(RoundFixed(1000.05214, 3), ShouldEqual, 1000.052)
			So(RoundFixed(math.Inf(1), 2), ShouldEqual, math.Inf(1))
			So(RoundFixed(math.Inf(-1), 2), ShouldEqual, math.Inf(-1))
		})

		Convey("Round", func() {
			So(Round(0.0171, 2), ShouldEqual, 0.017)
			So(Round(1234.567, 3), ShouldEqual, 1230.0)
			So(Round(math.Inf(1), 2), ShouldEqual, math.Inf(1))
			So(Round(math.Inf(-1), 2), ShouldEqual, math.Inf(-1))
		})

		Convey("RoundSlice", func() {
			So(RoundSlice([]float64{0.0012345, 12.345, 1234500.0}, 3),
				ShouldResemble, []float64{0.00123, 12.3, 1230000.0})
		})
	})
}
