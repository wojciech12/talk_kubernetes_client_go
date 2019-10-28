//

package main

import (
	"testing"

	v1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/fake"
)

func TestCheckPods(t *testing.T) {

	clientSet := fake.NewSimpleClientset(&v1.Pod{
		TypeMeta: metaV1.TypeMeta{
			Kind:       "Pod",
			APIVersion: "v1",
		},
		ObjectMeta: metaV1.ObjectMeta{
			Name:      "my-app",
			Namespace: "my-ns",
			Labels: map[string]string{
				"tag": "",
			},
		},
	})

	found := check(clientSet, "my-ns", "my-app")
	if !found {
		t.Fatal("Pod not found!")
	}

}
