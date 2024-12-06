/*
Copyright AppsCode Inc. and Contributors

Licensed under the AppsCode Community License 1.0.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package config

import (
	"bytes"
	"github.com/spf13/cobra"
	"io"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"kmodules.xyz/client-go/tools/parser"
	"os"
	"sigs.k8s.io/yaml"
)

func NewCmdCAPK() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "capk",
		Short:             "Configure CAPK config",
		DisableAutoGenTag: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			in, err := io.ReadAll(os.Stdin)
			if err != nil {
				return err
			}

			var out bytes.Buffer
			//	var foundCP bool

			err = parser.ProcessResources(in, func(ri parser.ResourceInfo) error {
				if ri.Object.GetKind() == "KubevirtMachineTemplate" {
					//foundCP = true

					if err := setBootstrapCheckStrategy(ri); err != nil {
						return err
					}

				}

				data, err := yaml.Marshal(ri.Object)
				if err != nil {
					return err
				}
				if out.Len() > 0 {
					out.WriteString("---\n")
				}
				_, err = out.Write(data)
				return err
			})
			if err != nil {
				return err
			}

			_, err = os.Stdout.Write(out.Bytes())
			return err
		},
	}

	return cmd
}

func setBootstrapCheckStrategy(ri parser.ResourceInfo) error {
	if err := unstructured.SetNestedField(ri.Object.UnstructuredContent(), "none", "spec", "template", "spec", "virtualMachineBootstrapCheck", "checkStrategy"); err != nil {
		return err
	}
	return nil
}