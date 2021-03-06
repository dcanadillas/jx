package kube_test

import (
	"bytes"
	"fmt"
	"github.com/jenkins-x/jx/pkg/kube"
	"github.com/jenkins-x/jx/pkg/testkube"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/fake"
	"strings"
	"testing"
)

func TestLogMasker(t *testing.T) {
	k8sObjects := []runtime.Object{
		testkube.CreateFakeGitSecret(),
	}
	ns := "jx"
	client := fake.NewSimpleClientset(k8sObjects...)

	hideValues := []string{
		"fakeuser",
		"fakepwd",
	}

	var buffer bytes.Buffer

	for i, hideValue := range hideValues {
		buffer.WriteString(fmt.Sprintf("%d: hide: %s\n", i+1, hideValue))
	}
	text := buffer.String()

	logMasker := &kube.LogMasker{}
	err := logMasker.LoadSecrets(client, ns)
	assert.NoError(t, err, "loading secrets in namespace %s", ns)

	actual := logMasker.MaskLog(text)

	t.Logf("created masked text: %s\n", actual)

	for _, hideValue := range hideValues {
		index := strings.Index(actual, hideValue)
		assert.True(t, index < 0, "found text %s at index %d in masked log: %s", hideValue, index, actual)
	}
}
