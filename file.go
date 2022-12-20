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
)

// WriteFile creates and writes a file with the given content.
func WriteFile(fileName, text string) error {
	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(text))
	return err
}

// FileExists returns true if a file with the given path exists, and it is
// indeed a file (not a directory).
func FileExists(fileName string) bool {
	info, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ReadFile returns the contents of the file as a string. If the file doesn't
// exist, or it is not a file, this method panics. It is a good practice to test
// that FileExists() is true, then call this method.
func ReadFile(fileName string) string {
	data, err := os.ReadFile(fileName)
	if err != nil {
		panic("cannot read file " + fileName + ": " + err.Error())
	}
	return string(data)
}
