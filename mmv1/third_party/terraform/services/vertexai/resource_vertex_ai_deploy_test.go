package vertexai_test

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-google/google/acctest"
	"github.com/hashicorp/terraform-provider-google/google/tpgresource"
	transport_tpg "github.com/hashicorp/terraform-provider-google/google/transport"
	"strings"
	"testing"
)

func TestAccVertexAIEndpointWithModelGardenDeployment_basic(t *testing.T) {
	t.Parallel()
	context := map[string]interface{}{"random_suffix": acctest.RandString(t, 10)}
	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEndpointWithModelGardenDeploymentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEndpointWithModelGardenDeployment_basic(context),
			},
		},
	})
}

func testAccVertexAIEndpointWithModelGardenDeployment_basic(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_endpoint_with_model_garden_deployment" "test" {
  publisher_model_name = "publishers/google/models/paligemma@paligemma-224-float32"
  location             = "us-central1"
  model_config {
    accept_eula =  true
  }
}
`, context)
}

func TestAccVertexAIEndpointWithModelGardenDeployment_withConfigs(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEndpointWithModelGardenDeploymentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEndpointWithModelGardenDeployment_withConfigs(context),
			},
		},
	})
}

func testAccVertexAIEndpointWithModelGardenDeployment_withConfigs(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_endpoint_with_model_garden_deployment" "test_with_configs" {
  publisher_model_name = "publishers/google/models/paligemma@paligemma-224-float32"
  location             = "us-central1"
  model_config {
    accept_eula =  true
  }
  deploy_config {
    dedicated_resources {
      machine_spec {
        machine_type      = "g2-standard-16"
        accelerator_type  = "NVIDIA_L4"
        accelerator_count = 1
      }
      min_replica_count = 1
    }
  }
}
`, context)
}

func TestAccVertexAIEndpointWithModelGardenDeployment_huggingfaceModel(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEndpointWithModelGardenDeploymentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEndpointWithModelGardenDeployment_huggingfaceModel(context),
			},
		},
	})
}

func testAccVertexAIEndpointWithModelGardenDeployment_huggingfaceModel(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_endpoint_with_model_garden_deployment" "deploy" {
  hugging_face_model_id = "Qwen/Qwen3-0.6B"
  location             = "us-central1"
  model_config {
    accept_eula = true
  }
}
`, context)
}

func TestAccVertexAIEndpointWithModelGardenDeployment_multipleModelsInSequence(t *testing.T) {
	t.Parallel()

	context := map[string]interface{}{
		"random_suffix": acctest.RandString(t, 10),
	}

	acctest.VcrTest(t, resource.TestCase{
		PreCheck:                 func() { acctest.AccTestPreCheck(t) },
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories(t),
		CheckDestroy:             testAccCheckVertexAIEndpointWithModelGardenDeploymentDestroyProducer(t),
		Steps: []resource.TestStep{
			{
				Config: testAccVertexAIEndpointWithModelGardenDeployment_multipleModelsInSequence(context),
			},
		},
	})
}

func testAccVertexAIEndpointWithModelGardenDeployment_multipleModelsInSequence(context map[string]interface{}) string {
	return acctest.Nprintf(`
resource "google_vertex_ai_endpoint_with_model_garden_deployment" "deploy-gemma-1_1-2b-it" {
  publisher_model_name = "publishers/google/models/gemma@gemma-1.1-2b-it"
  location             = "us-central1"
  model_config {
    accept_eula = true
  }
  deploy_config {
    dedicated_resources {
      machine_spec {
        machine_type      = "g2-standard-12"
        accelerator_type  = "NVIDIA_L4"
        accelerator_count = 1
      }
      min_replica_count = 1
    }
  }
}

resource "google_vertex_ai_endpoint_with_model_garden_deployment" "deploy-qwen3-0_6b" {
  hugging_face_model_id = "Qwen/Qwen3-0.6B"
  location             = "us-central1"
  model_config {
    accept_eula = true
  }
  deploy_config {
    dedicated_resources {
      machine_spec {
        machine_type      = "g2-standard-12"
        accelerator_type  = "NVIDIA_L4"
        accelerator_count = 1
      }
      min_replica_count = 1
    }
  }
  depends_on = [ google_vertex_ai_endpoint_with_model_garden_deployment.deploy-gemma-1_1-2b-it ]
}

resource "google_vertex_ai_endpoint_with_model_garden_deployment" "deploy-llama-3_2-1b" {
  publisher_model_name = "publishers/meta/models/llama3-2@llama-3.2-1b"
  location             = "us-central1"
  model_config {
    accept_eula = true
  }
  deploy_config {
    dedicated_resources {
      machine_spec {
        machine_type      = "g2-standard-12"
        accelerator_type  = "NVIDIA_L4"
        accelerator_count = 1
      }
      min_replica_count = 1
    }
  }
  depends_on = [ google_vertex_ai_endpoint_with_model_garden_deployment.deploy-qwen3-0_6b ]
}
`, context)
}

func testAccCheckVertexAIEndpointWithModelGardenDeploymentDestroyProducer(t *testing.T) func(s *terraform.State) error {
	return func(s *terraform.State) error {
		for name, rs := range s.RootModule().Resources {
			if rs.Type != "google_vertex_ai_endpoint_with_model_garden_deployment" {
				continue
			}
			if strings.HasPrefix(name, "data.") {
				continue
			}

			config := acctest.GoogleProviderConfig(t)

			url, err := tpgresource.ReplaceVarsForTest(config, rs, "{{VertexAIBasePath}}projects/{{project}}/locations/{{location}}/endpoints/{{endpoint}}")
			if err != nil {
				return err
			}

			billingProject := ""

			if config.BillingProject != "" {
				billingProject = config.BillingProject
			}

			_, err = transport_tpg.SendRequest(transport_tpg.SendRequestOptions{
				Config:    config,
				Method:    "GET",
				Project:   billingProject,
				RawURL:    url,
				UserAgent: config.UserAgent,
			})
			if err == nil {
				return fmt.Errorf("VertexAIEndpointWithModelGardenDeployment still exists at %s", url)
			}
		}

		return nil
	}
}
