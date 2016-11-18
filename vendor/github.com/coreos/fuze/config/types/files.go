// Copyright 2016 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package types

type File struct {
	Filesystem string       `yaml:"filesystem"`
	Path       string       `yaml:"path"`
	Mode       int          `yaml:"mode"`
	Contents   FileContents `yaml:"contents"`
	User       FileUser     `yaml:"user"`
	Group      FileGroup    `yaml:"group"`
}

type FileContents struct {
	Remote Remote `yaml:"remote"`
	Inline string `yaml:"inline"`
}

type Remote struct {
	Url          string       `yaml:"url"`
	Compression  string       `yaml:"compression"`
	Verification Verification `yaml:"verification"`
}

type FileUser struct {
	Id int `yaml:"id"`
}

type FileGroup struct {
	Id int `yaml:"id"`
}
