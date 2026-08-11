package gcpcredential

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

type mockTransport struct {
	roundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.roundTripFunc(req)
}

func TestProvide_WorkloadIdentity(t *testing.T) {
	validToken := "dummyHeader.eyJpc3MiOiAiaHR0cHM6Ly9jb250YWluZXIuZ29vZ2xlYXBpcy5jb20vdjEvcHJvamVjdHMvbXktcHJvamVjdC9sb2NhdGlvbnMvdXMtY2VudHJhbDEvY2x1c3RlcnMvbXktY2x1c3RlciJ9.dummySignature"
	tests := []struct {
		name                string
		identityProvider    string
		serviceAccountToken string
		metadataResponses   map[string]string
		stsResponse         string
		wantAudience        string
		expectedToken       string
	}{
		{
			name:                "Direct Access Mode (Success)",
			identityProvider:    "//iam.googleapis.com/projects/my-project-number/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
			serviceAccountToken: validToken,
			stsResponse:         `{"access_token": "federated-token-xyz", "expires_in": 3600, "token_type": "Bearer"}`,
			wantAudience:        "//iam.googleapis.com/projects/my-project-number/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
			expectedToken:       "federated-token-xyz",
		},
		{
			name:                "Fallback to Node SA (No Identity Provider, bypass STS)",
			serviceAccountToken: validToken,
			metadataResponses: map[string]string{
				"project/project-id":                      "my-project",
				"instance/service-accounts/default/token": `{"access_token": "node-sa-token", "expires_in": 3600}`,
				"instance/service-accounts/default/email": "node-sa@project.gserviceaccount.com",
				"instance/service-accounts/":              "default/\n",
			},
			stsResponse:   `{"access_token": "should-not-be-called", "expires_in": 3600, "token_type": "Bearer"}`,
			expectedToken: "node-sa-token",
		},
		{
			name:                "Fail Fast - Identity Provider Configured, STS Fails",
			identityProvider:    "//iam.googleapis.com/projects/my-project-number/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
			serviceAccountToken: validToken,
			metadataResponses: map[string]string{
				"instance/service-accounts/default/token": `{"access_token": "node-sa-token", "expires_in": 3600}`,
			},
			stsResponse:   "error",
			wantAudience:  "//iam.googleapis.com/projects/my-project-number/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
			expectedToken: "",
		},
		{
			name:                "Fail Fast - Identity Provider Configured, Token is Empty",
			identityProvider:    "//iam.googleapis.com/projects/my-project-number/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
			serviceAccountToken: "",
			metadataResponses: map[string]string{
				"instance/service-accounts/default/token": `{"access_token": "node-sa-token", "expires_in": 3600}`,
			},
			stsResponse:   `{"access_token": "should-not-be-called", "expires_in": 3600, "token_type": "Bearer"}`,
			expectedToken: "",
		},
		{
			name:                "Direct Access Mode - Configured Identity Provider (Success)",
			identityProvider:    "https://custom-provider.com",
			serviceAccountToken: validToken,
			stsResponse:         `{"access_token": "federated-token-custom", "expires_in": 3600, "token_type": "Bearer"}`,
			wantAudience:        "https://custom-provider.com",
			expectedToken:       "federated-token-custom",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mTransport := &mockTransport{
				roundTripFunc: func(req *http.Request) (*http.Response, error) {
					// Handle Metadata Server requests
					if req.URL.Host == "metadata.google.internal." {
						path := strings.TrimPrefix(req.URL.Path, "/computeMetadata/v1/")
						if resp, ok := tc.metadataResponses[path]; ok {
							return &http.Response{
								StatusCode: http.StatusOK,
								Body:       io.NopCloser(bytes.NewBufferString(resp)),
								Header:     make(http.Header),
							}, nil
						}
						return &http.Response{
							StatusCode: http.StatusNotFound,
							Body:       io.NopCloser(bytes.NewBufferString("")),
						}, nil
					}

					// Handle STS Requests
					if req.URL.Host == "sts.googleapis.com" && req.URL.Path == "/v1/token" {
						if tc.stsResponse == "error" {
							return &http.Response{
								StatusCode: http.StatusBadRequest,
								Body:       io.NopCloser(bytes.NewBufferString(`{"error": "invalid_grant"}`)),
							}, nil
						}

						// Verify audience if specified
						if tc.wantAudience != "" {
							bodyBytes, err := io.ReadAll(req.Body)
							if err != nil {
								t.Errorf("failed to read STS request body: %v", err)
							}
							var stsReq stsTokenExchangeRequest
							if err := json.Unmarshal(bodyBytes, &stsReq); err != nil {
								t.Errorf("failed to unmarshal STS request body: %v", err)
							}
							if stsReq.Audience != tc.wantAudience {
								t.Errorf("unexpected STS audience: got=%q, want=%q", stsReq.Audience, tc.wantAudience)
							}
						}

						return &http.Response{
							StatusCode: http.StatusOK,
							Body:       io.NopCloser(bytes.NewBufferString(tc.stsResponse)),
						}, nil
					}

					return &http.Response{
						StatusCode: http.StatusNotFound,
						Body:       io.NopCloser(bytes.NewBufferString("")),
					}, nil
				},
			}

			httpClient := &http.Client{
				Transport: mTransport,
			}

			provider := &ContainerRegistryProvider{
				MetadataProvider: MetadataProvider{
					Client: httpClient,
				},
				UseRegistryFromImage: true,
			}
			provider.KSAToken = tc.serviceAccountToken
			provider.IdentityProvider = tc.identityProvider

			cfg := provider.Provide("us-central1-docker.pkg.dev/my-project/my-repo/my-image:latest")

			if tc.expectedToken == "" {
				for _, entry := range cfg {
					if entry.Password != "" {
						t.Errorf("expected no token, got: %s", entry.Password)
					}
				}
			} else {
				found := false
				for _, entry := range cfg {
					if entry.Password == tc.expectedToken {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected token %q not found in config: %+v", tc.expectedToken, cfg)
				}
			}
		})
	}
}
