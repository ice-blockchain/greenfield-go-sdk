package client

import (
	"net/url"
	"testing"

	"github.com/bnb-chain/greenfield-go-sdk/types"
)

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatalf("invalid test URL %q: %v", raw, err)
	}
	return u
}

func TestGenerateURL(t *testing.T) {
	t.Parallel()
	c := &Client{}

	tests := []struct {
		name          string
		bucketName    string
		objectName    string
		relativePath  string
		queryValues   url.Values
		adminInfo     AdminAPIInfo
		endpoint      string
		isVirtualHost bool
		want          string
		wantErr       bool
	}{
		{
			name:       "basic path-style bucket only",
			endpoint:   "https://sp.example.com",
			bucketName: "mybucket",
			want:       "https://sp.example.com/mybucket/",
		},
		{
			name:       "path-style bucket and object",
			endpoint:   "https://sp.example.com",
			bucketName: "mybucket",
			objectName: "myobject",
			want:       "https://sp.example.com/mybucket/myobject/",
		},
		{
			name:          "virtual host bucket only",
			endpoint:      "https://sp.example.com",
			bucketName:    "mybucket",
			isVirtualHost: true,
			want:          "https://mybucket.sp.example.com/",
		},
		{
			name:          "virtual host bucket and object",
			endpoint:      "https://sp.example.com",
			bucketName:    "mybucket",
			objectName:    "myobject",
			isVirtualHost: true,
			want:          "https://mybucket.sp.example.com/myobject/",
		},
		{
			name:     "no bucket or object",
			endpoint: "https://sp.example.com",
			want:     "https://sp.example.com/",
		},
		{
			name:       "endpoint with path component",
			endpoint:   "https://sp.example.com/base/path",
			bucketName: "mybucket",
			objectName: "myobject",
			want:       "https://sp.example.com/base/path/mybucket/myobject/",
		},
		{
			name:       "object name with spaces (no double encoding)",
			endpoint:   "https://sp.example.com",
			bucketName: "mybucket",
			objectName: "my object",
			want:       "https://sp.example.com/mybucket/my%20object/",
		},
		{
			name:       "object name with special chars",
			endpoint:   "https://sp.example.com",
			bucketName: "mybucket",
			objectName: "path/to/my file.txt",
			want:       "https://sp.example.com/mybucket/path/to/my%20file.txt/",
		},
		{
			name:       "object name with unicode chars",
			endpoint:   "https://sp.example.com",
			bucketName: "mybucket",
			objectName: "日本語.txt",
			want:       "https://sp.example.com/mybucket/%E6%97%A5%E6%9C%AC%E8%AA%9E.txt/",
		},
		{
			name:         "with relative path",
			endpoint:     "https://sp.example.com",
			bucketName:   "mybucket",
			relativePath: "extra/segment",
			want:         "https://sp.example.com/mybucket/extra/segment",
		},
		{
			name:        "with query values",
			endpoint:    "https://sp.example.com",
			bucketName:  "mybucket",
			queryValues: url.Values{"key": {"value"}},
			want:        "https://sp.example.com/mybucket/?key=value",
		},
		{
			name:       "strip default https port 443",
			endpoint:   "https://sp.example.com:443",
			bucketName: "mybucket",
			want:       "https://sp.example.com/mybucket/",
		},
		{
			name:       "strip default http port 80",
			endpoint:   "http://sp.example.com:80",
			bucketName: "mybucket",
			want:       "http://sp.example.com/mybucket/",
		},
		{
			name:       "keep non-default port",
			endpoint:   "https://sp.example.com:9090",
			bucketName: "mybucket",
			want:       "https://sp.example.com:9090/mybucket/",
		},
		{
			name:     "admin API v1",
			endpoint: "https://sp.example.com",
			adminInfo: AdminAPIInfo{
				isAdminAPI:   true,
				adminVersion: types.AdminV1Version,
			},
			want: "https://sp.example.com" + types.AdminURLPrefix + types.AdminURLV1Version + "/",
		},
		{
			name:     "admin API v2",
			endpoint: "https://sp.example.com",
			adminInfo: AdminAPIInfo{
				isAdminAPI:   true,
				adminVersion: types.AdminV2Version,
			},
			want: "https://sp.example.com" + types.AdminURLPrefix + types.AdminURLV2Version + "/",
		},
		{
			name:     "admin API invalid version",
			endpoint: "https://sp.example.com",
			adminInfo: AdminAPIInfo{
				isAdminAPI:   true,
				adminVersion: 999,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			endpoint := mustParseURL(t, tt.endpoint)
			got, err := c.generateURL(tt.bucketName, tt.objectName, tt.relativePath,
				tt.queryValues, tt.adminInfo, endpoint, tt.isVirtualHost)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got.String() != tt.want {
				t.Errorf("got  %q\nwant %q", got.String(), tt.want)
			}
		})
	}
}
