/*
Copyright 2021 The Dapr Authors
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

package kubernetes

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dapr/kit/logger"
)

func TestGetNamespace(t *testing.T) {
	t.Run("has namespace metadata", func(t *testing.T) {
		store := kubernetesSecretStore{logger: logger.NewLogger("test")}
		namespace := "a"

		ns, err := store.getNamespaceFromMetadata(map[string]string{"namespace": namespace})
		require.NoError(t, err)
		assert.Equal(t, namespace, ns)
	})

	t.Run("has namespace env", func(t *testing.T) {
		store := kubernetesSecretStore{logger: logger.NewLogger("test")}
		t.Setenv("NAMESPACE", "b")

		ns, err := store.getNamespaceFromMetadata(map[string]string{})
		require.NoError(t, err)
		assert.Equal(t, "b", ns)
	})

	t.Run("no namespace", func(t *testing.T) {
		store := kubernetesSecretStore{logger: logger.NewLogger("test")}
		t.Setenv("NAMESPACE", "")
		_, err := store.getNamespaceFromMetadata(map[string]string{})

		require.Error(t, err)
		require.ErrorContains(t, err, "namespace is missing")
	})

	t.Run("has default namespace", func(t *testing.T) {
		store := kubernetesSecretStore{
			logger: logger.NewLogger("test"),
			md: kubernetesMetadata{
				DefaultNamespace: "c",
			},
		}

		ns, err := store.getNamespaceFromMetadata(map[string]string{})
		require.NoError(t, err)
		assert.Equal(t, "c", ns)
	})
}

func TestGetFeatures(t *testing.T) {
	s := kubernetesSecretStore{logger: logger.NewLogger("test")}
	// Yes, we are skipping initialization as feature retrieval doesn't depend on it.
	t.Run("no features are advertised", func(t *testing.T) {
		f := s.Features()
		assert.Empty(t, f)
	})
}
