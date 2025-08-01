package acctest_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"testing"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
)

func TestNewVcrMatcherFunc_canDetectMatches(t *testing.T) {

	// Everything should be determined as a match in this test
	cases := map[string]struct {
		httpRequest     requestDescription
		cassetteRequest requestDescription
	}{
		"matches POST requests with empty bodies": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{}",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{}",
			},
		},
		"matches POST requests with exact matching bodies": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value\"}",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value\"}",
			},
		},
		"matches POST requests with matching but re-ordered bodies, but only if Content-Type contains 'application/json'": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: "{\"field1\":\"value1\",\"field2\":\"value2\"}", // 1 before 2
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"Content-Type": "application/json",
				},
				body: "{\"field2\":\"value2\",\"field1\":\"value1\"}", // 2 before 1
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Make matcher
			ctx := context.Background()
			req := prepareHttpRequest(tc.httpRequest)
			cassetteReq := prepareCassetteRequest(tc.cassetteRequest)
			matcher := acctest.NewVcrMatcherFunc(ctx)

			// Act - use matcher
			matchDetected := matcher(req, cassetteReq)

			// Assert match
			if !matchDetected {
				t.Fatalf("expected matcher to match the requests")
			}
		})
	}
}

func TestNewVcrMatcherFunc_canDetectMismatches(t *testing.T) {

	// All these cases are expected to end with no match detected
	cases := map[string]struct {
		httpRequest     requestDescription
		cassetteRequest requestDescription
	}{
		"requests using different schemes": {
			httpRequest: requestDescription{
				scheme: "http",
				method: "GET",
				host:   "example.com",
				path:   "foobar",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar",
			},
		},
		"requests using different hosts": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "google.com",
				path:   "foobar",
			},
		},
		"requests using different paths": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar1",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar2",
			},
		},
		"requests with different methods": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{}",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "PUT",
				host:   "example.com",
				path:   "foobar",
				body:   "{}",
			},
		},
		"POST requests with different bodies": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value is ABCDEFG\"}",
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value is MNLOP\"}",
			},
		},
		"POST requests with matching but re-ordered bodies aren't matching if Content-Type header is not 'application/json'": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"Content-Type": "foobar",
				},
				body: "{\"field1\":\"value1\",\"field2\":\"value2\"}", // 1 before 2
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"Content-Type": "foobar",
				},
				body: "{\"field2\":\"value2\",\"field1\":\"value1\"}", // 2 before 1
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Make matcher
			ctx := context.Background()
			req := prepareHttpRequest(tc.httpRequest)
			cassetteReq := prepareCassetteRequest(tc.cassetteRequest)
			matcher := acctest.NewVcrMatcherFunc(ctx)

			// Act - use matcher
			matchDetected := matcher(req, cassetteReq)

			// Assert match
			if matchDetected {
				t.Fatalf("expected matcher to not match the requests")
			}
		})
	}
}

// Currently there is no code to actively force the matcher to overlook differing User-Agent values.
// It isn't checked at any point in the matcher logic.
func TestNewVcrMatcherFunc_ignoresDifferentUserAgents(t *testing.T) {

	cases := map[string]struct {
		httpRequest     requestDescription
		cassetteRequest requestDescription
	}{
		"GET requests with different useragents are matched": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"User-Agent": "user-agent-HTTP",
				},
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "GET",
				host:   "example.com",
				path:   "foobar",
				headers: map[string]string{
					"User-Agent": "user-agent-CASSETTE",
				},
			},
		},
		"POST requests with identical bodies and different useragents are matched": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value\"}",
				headers: map[string]string{
					"User-Agent": "user-agent-HTTP",
				},
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field\":\"value\"}",
				headers: map[string]string{
					"User-Agent": "user-agent-CASSETTE",
				},
			},
		},
		"POST requests with reordered but matching bodies and different useragents are matched if Content-Type contains 'application/json'": {
			httpRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field1\":\"value1\",\"field2\":\"value2\"}",
				headers: map[string]string{
					"User-Agent":   "user-agent-HTTP",
					"Content-Type": "application/json",
				},
			},
			cassetteRequest: requestDescription{
				scheme: "https",
				method: "POST",
				host:   "example.com",
				path:   "foobar",
				body:   "{\"field2\":\"value2\",\"field1\":\"value1\"}",
				headers: map[string]string{
					"User-Agent":   "user-agent-CASSETTE",
					"Content-Type": "application/json",
				},
			},
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			// Make matcher
			ctx := context.Background()
			req := prepareHttpRequest(tc.httpRequest)
			cassetteReq := prepareCassetteRequest(tc.cassetteRequest)
			matcher := acctest.NewVcrMatcherFunc(ctx)

			// Act - use matcher
			matchDetected := matcher(req, cassetteReq)

			// Assert match
			if !matchDetected {
				t.Fatalf("expected matcher to match the requests")
			}
		})
	}
}

type requestDescription struct {
	scheme  string
	method  string
	host    string
	path    string
	body    string
	headers map[string]string
}

func prepareHttpRequest(d requestDescription) *http.Request {
	url := &url.URL{
		Scheme: d.scheme,
		Host:   d.host,
		Path:   d.path,
	}

	req := &http.Request{
		Method: d.method,
		URL:    url,
	}

	// Conditionally set a body
	if d.body != "" {
		body := io.NopCloser(bytes.NewBufferString(d.body))
		req.Body = body
	}
	// Conditionally set headers
	if len(d.headers) > 0 {
		req.Header = http.Header{}
		for k, v := range d.headers {
			req.Header.Set(k, v)
		}
	}

	return req
}

func prepareCassetteRequest(d requestDescription) cassette.Request {
	fullUrl := fmt.Sprintf("%s://%s/%s", d.scheme, d.host, d.path)

	req := cassette.Request{
		Method: d.method,
		URL:    fullUrl,
	}

	// Conditionally set a body
	if d.body != "" {
		req.Body = d.body
	}
	// Conditionally set headers
	if len(d.headers) > 0 {
		req.Headers = http.Header{}
		for k, v := range d.headers {
			req.Headers.Add(k, v)
		}
	}

	return req
}

func TestReformConfigWithProvider(t *testing.T) {

	type testCase struct {
		name             string
		initialConfig    string
		providerToInsert string
		expectedConfig   string
	}

	cases := map[string]testCase{
		"replaces_google_beta_with_local": {
			name: "Replaces 'google-beta' provider with 'google-local'",
			initialConfig: `resource "google_new_resource" {
      provider = google-beta
}`,
			providerToInsert: "google-local",
			expectedConfig: `resource "google_new_resource" {
      provider = google-local
}`,
		},
		"inserts_local_provider_into_empty_config": {
			name: "Inserts 'google-local' provider when no provider block exists",
			initialConfig: `resource "google_alloydb_cluster" "default" {
    location   = "us-central1"
    network_config {
        network = google_compute_network.default.id
    }
}`,
			providerToInsert: "google-local",
			expectedConfig: `resource "google_alloydb_cluster" "default" {
  provider = google-local

    location   = "us-central1"
    network_config {
        network = google_compute_network.default.id
    }
}`,
		},
		"no_change_if_target_provider_already_present": {
			name: "Does not change config if target provider is already present",
			initialConfig: `resource "google_new_resource" {
      provider = google-local
}`,
			providerToInsert: "google-local",
			expectedConfig: `resource "google_new_resource" {
      provider = google-local
}`,
		},
		"inserts_provider_with_other_attributes": {
			name: "Inserts provider into a resource block with other attributes but no existing provider",
			initialConfig: `resource "google_compute_instance" "test" {
  name         = "test-instance"
  machine_type = "e2-medium"
}`,
			providerToInsert: "google-local",
			expectedConfig: `resource "google_compute_instance" "test" {
  provider = google-local

  name         = "test-instance"
  machine_type = "e2-medium"
}`,
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			newConfig := acctest.ReformConfigWithProvider(tc.initialConfig, tc.providerToInsert)

			if newConfig != tc.expectedConfig {
				t.Fatalf("Test Case: %s\nExpected config to be reformatted to:\n%q\nbut got:\n%q", tc.name, tc.expectedConfig, newConfig)
			}
			t.Logf("Test Case: %s\nReformed config:\n%s", tc.name, newConfig)
		})
	}
}

func TestInsertDiffSteps(t *testing.T) {

	var dummyCase = resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		Steps: []resource.TestStep{
			{
				Config: `resource "google_new_resource" "original" {
                    provider = google-beta
                }`,
			},
			{
				Config: `resource "google_new_resource" "original" {
                    provider = google-beta
                }`,
			},
			{
				ResourceName:            "google_pubsub_subscription.example",
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"topic"},
			},
			{
				Config: `resource "google_example_widget" "foo" {
					name = "dummy"
					provider = google-beta
				}`,
				Check: resource.ComposeTestCheckFunc(
					func(*terraform.State) error { return nil },
				),
			},
			{
				Config: `provider = "google-local"
						// ... configuration that is expected to cause an error
					`,
				ExpectError: regexp.MustCompile(`"restore_continuous_backup_source": conflicts with restore_backup_source`),
			},
		},
	}
	temp_file, err := os.CreateTemp("", "release_diff_test_output_*.log")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	dummyCase = acctest.InsertDiffSteps(dummyCase, temp_file, "google-beta", "google-local")

	// Expected steps after InsertDiffSteps runs.
	// A "diff" step (using 'google-local') is added for each original step containing a Config field,
	// unless the step has ExpectError set.
	var expectedSteps = []resource.TestStep{
		{
			Config: `resource "google_new_resource" "original" {
                    provider = google-beta
                }`,
		},
		{
			Config: `resource "google_new_resource" "original" {
                    provider = google-local
                }`,
			ExpectNonEmptyPlan: false,
			PlanOnly:           true,
		},
		{
			Config: `resource "google_new_resource" "original" {
                    provider = google-beta
                }`,
		},
		{
			Config: `resource "google_new_resource" "original" {
                    provider = google-local
                }`,
			ExpectNonEmptyPlan: false,
			PlanOnly:           true,
		},
		{
			ResourceName: "google_pubsub_subscription.example", // No config, so no diff step added
		},
		{
			Config: `resource "google_example_widget" "foo" {
					name = "dummy"
					provider = google-beta
				}`,
			Check: resource.ComposeTestCheckFunc(
				func(*terraform.State) error { return nil },
			),
		},
		{
			Config: `resource "google_example_widget" "foo" {
					name = "dummy"
					provider = google-local
				}`,
			Check: resource.ComposeTestCheckFunc(
				func(*terraform.State) error { return nil },
			),
			ExpectNonEmptyPlan: false,
			PlanOnly:           true,
		},
		{
			Config: `provider = "google-local"
						// ... configuration that is expected to cause an error
					`, // expect error means we don't do a second step
		},
	}

	if len(dummyCase.Steps) != len(expectedSteps) {
		t.Fatalf("Expected %d steps, but got %d", len(expectedSteps), len(dummyCase.Steps))
	}

	for i, step := range dummyCase.Steps {
		if step.Config != expectedSteps[i].Config {
			t.Fatalf("Expected step %d config to be:\n%q\nbut got:\n%q", i, expectedSteps[i].Config, step.Config)
		}
		if step.PlanOnly != expectedSteps[i].PlanOnly {
			t.Fatalf("Expected step %d to have PlanOnly set to %v, but got %v", i, expectedSteps[i].PlanOnly, step.PlanOnly)
		}
	}

	defer os.Remove(temp_file.Name())
}
