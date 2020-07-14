package android

import "testing"

func TestGitHubContent_Decode(t *testing.T) {
	content := GitHubContent{
		Encoding: "base64",
		Name:     "build.gradle",
		Content:  "YXBwbHkgcGx1Z2luOiAnY29tLmFuZHJvaWQuYXBwbGljYXRpb24nCmFwcGx5\nIHBsdWdpbjogJ2tvdGxpbi1hbmRyb2lkJwphcHBseSBwbHVnaW46ICdrb3Rs\naW4tYW5kcm9pZC1leHRlbnNpb25zJwphcHBseSBwbHVnaW46ICdrb3RsaW4t\n",
	}

	decoded, err := content.Decode()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Decoced content: %s", decoded)
}
