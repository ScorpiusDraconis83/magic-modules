---
subcategory: "Parameter Manager"
description: |-
  Get information about an Parameter Manager Parameter Version
---

# google_parameter_manager_parameter_version

Get the value and metadata from a Parameter Manager Parameter version. For more information see the [official documentation](https://cloud.google.com/secret-manager/parameter-manager/docs/overview)  and [API](https://cloud.google.com/secret-manager/parameter-manager/docs/reference/rest/v1/projects.locations.parameters.versions).

## Example Usage

```hcl
data "google_parameter_manager_parameter_version" "basic" {
  parameter            = "test-parameter"
  parameter_version_id = "test-parameter-version"
}
```

## Argument Reference

The following arguments are supported:

* `project` - (Optional) The project for retrieving the Parameter Version. If it's not specified, 
    the provider project will be used.

* `parameter` - (Required) The parameter for obtaining the Parameter Version.
    This can be either the reference of the parameter as in `projects/{{project}}/locations/global/parameters/{{parameter_id}}` or only the name of the parameter as in `{{parameter_id}}`.

* `parameter_version_id` - (Required) The version of the parameter to get.

## Attributes Reference

The following attributes are exported:

* `parameter_data` - The parameter data.

* `name` - The resource name of the ParameterVersion. Format:
  `projects/{{project}}/locations/global/parameters/{{parameter_id}}/versions/{{parameter_version_id}}`

* `create_time` - The time at which the Parameter Version was created.

* `update_time` - The time at which the Parameter Version was last updated.

* `disabled` -  The current state of the Parameter Version.

* `kms_key_version` - The resource name of the Cloud KMS CryptoKeyVersion used to decrypt parameter version payload. Format `projects/{{project}}/locations/global/keyRings/{{key_ring}}/cryptoKeys/{{crypto_key}}/cryptoKeyVersions/{{crypto_key_version}}`
