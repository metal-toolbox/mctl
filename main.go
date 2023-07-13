/*
Copyright © 2022 Equinix Metal <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"github.com/metal-toolbox/mctl/cmd"
	_ "github.com/metal-toolbox/mctl/cmd/create"
	_ "github.com/metal-toolbox/mctl/cmd/delete"
	_ "github.com/metal-toolbox/mctl/cmd/edit"
	_ "github.com/metal-toolbox/mctl/cmd/generate"
	_ "github.com/metal-toolbox/mctl/cmd/get"
	_ "github.com/metal-toolbox/mctl/cmd/install"
	_ "github.com/metal-toolbox/mctl/cmd/list"
)

func main() {
	cmd.Execute()
}
