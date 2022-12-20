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
	"os"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFile(t *testing.T) {
	t.Parallel()

	Convey("File methods work", t, func() {
		tmpdir, tmpdirErr := os.MkdirTemp("", "testmain")
		defer os.RemoveAll(tmpdir)

		So(tmpdirErr, ShouldBeNil)

		fileName := filepath.Join(tmpdir, "test.txt")

		Convey("Write, check for existence and read", func() {
			So(FileExists(fileName), ShouldBeFalse)
			text := "Time flies like an arrow. To them, it's tasty."
			So(WriteFile(fileName, text), ShouldBeNil)
			So(FileExists(fileName), ShouldBeTrue)
			So(ReadFile(fileName), ShouldEqual, text)
		})

	})

}
