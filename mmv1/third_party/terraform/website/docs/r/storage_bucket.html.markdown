---
subcategory: "Cloud Storage"
description: |-
  Creates a new bucket in Google Cloud Storage.
---

# google_storage_bucket

Creates a new bucket in Google cloud storage service (GCS).
Once a bucket has been created, its location can't be changed.

For more information see
[the official documentation](https://cloud.google.com/storage/docs/overview)
and
[API](https://cloud.google.com/storage/docs/json_api/v1/buckets).

**Note**: If the project id is not set on the resource or in the provider block it will be dynamically
determined which will require enabling the compute api.


## Example Usage - creating a private bucket in standard storage, in the EU region. Bucket configured as static website and CORS configurations

```hcl
resource "google_storage_bucket" "static-site" {
  name          = "image-store.com"
  location      = "EU"
  force_destroy = true

  uniform_bucket_level_access = true

  website {
    main_page_suffix = "index.html"
    not_found_page   = "404.html"
  }
  cors {
    origin          = ["http://image-store.com"]
    method          = ["GET", "HEAD", "PUT", "POST", "DELETE"]
    response_header = ["*"]
    max_age_seconds = 3600
  }
}
```

## Example Usage - Life cycle settings for storage bucket objects

```hcl
resource "google_storage_bucket" "auto-expire" {
  name          = "auto-expiring-bucket"
  location      = "US"
  force_destroy = true

  lifecycle_rule {
    condition {
      age = 3
    }
    action {
      type = "Delete"
    }
  }

  lifecycle_rule {
    condition {
      age = 1
    }
    action {
      type = "AbortIncompleteMultipartUpload"
    }
  }
}
```

## Example Usage - Life cycle settings for storage bucket objects with `send_age_if_zero` disabled
When creating a life cycle condition that does not also include an `age` field, a default `age` of 0 will be set. Set the `send_age_if_zero` flag to `false` to prevent this and avoid any potentially unintended interactions.

```hcl
resource "google_storage_bucket" "no-age-enabled" {
  provider = google-beta
  name          = "no-age-enabled-bucket"
  location      = "US"
  force_destroy = true

  lifecycle_rule {
    action {
      type = "Delete"
    }
    condition {
      days_since_noncurrent_time = 3
      send_age_if_zero = false
    }
  }
}
```

## Example Usage - Enabling public access prevention

```hcl
resource "google_storage_bucket" "no-public-access" {
  name          = "no-public-access-bucket"
  location      = "US"
  force_destroy = true

  public_access_prevention = "enforced"
}
```

## Example Usage - Enabling hierarchical namespace

```hcl
resource "google_storage_bucket" "hns-enabled" {
  name          = "hns-enabled-bucket"
  location      = "US"
  force_destroy = true

  hierarchical_namespace {
    enabled = true
  }
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the bucket. Bucket names must be in lowercase and no more than 63 characters long. You can find the complete list of bucket naming rules [here](https://cloud.google.com/storage/docs/buckets#naming).

* `location` - (Required) The [GCS location](https://cloud.google.com/storage/docs/bucket-locations).

- - -

* `force_destroy` - (Optional, Default: false) When deleting a bucket, this
    boolean option will delete all contained objects. If you try to delete a
    bucket that contains objects, Terraform will fail that run.

* `project` - (Optional) The ID of the project in which the resource belongs. If it
    is not provided, the provider project is used.

* `storage_class` - (Optional, Default: 'STANDARD') The [Storage Class](https://cloud.google.com/storage/docs/storage-classes) of the new bucket. Supported values include: `STANDARD`, `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`, `ARCHIVE`.

* `autoclass` - (Optional) The bucket's [Autoclass](https://cloud.google.com/storage/docs/autoclass) configuration.  Structure is [documented below](#nested_autoclass).

* `lifecycle_rule` - (Optional) The bucket's [Lifecycle Rules](https://cloud.google.com/storage/docs/lifecycle#configuration) configuration. Multiple blocks of this type are permitted. Structure is [documented below](#nested_lifecycle_rule).

* `versioning` - (Optional) The bucket's [Versioning](https://cloud.google.com/storage/docs/object-versioning) configuration.  Structure is [documented below](#nested_versioning).

* `website` - (Optional) Configuration if the bucket acts as a website. Structure is [documented below](#nested_website).

* `cors` - (Optional) The bucket's [Cross-Origin Resource Sharing (CORS)](https://www.w3.org/TR/cors/) configuration. Multiple blocks of this type are permitted. Structure is [documented below](#nested_cors).

* `default_event_based_hold` - (Optional) Whether or not to automatically apply an eventBasedHold to new objects added to the bucket.

* `retention_policy` - (Optional) Configuration of the bucket's data retention policy for how long objects in the bucket should be retained. Structure is [documented below](#nested_retention_policy).

* `labels` - (Optional) A map of key/value label pairs to assign to the bucket.

* `logging` - (Optional) The bucket's [Access & Storage Logs](https://cloud.google.com/storage/docs/access-logs) configuration. Structure is [documented below](#nested_logging).

* `encryption` - (Optional) The bucket's encryption configuration. Structure is [documented below](#nested_encryption).

* `enable_object_retention` - (Optional, Default: false) Enables [object retention](https://cloud.google.com/storage/docs/object-lock) on a storage bucket.


* `requester_pays` - (Optional, Default: false) Enables [Requester Pays](https://cloud.google.com/storage/docs/requester-pays) on a storage bucket.

* `rpo` - (Optional) The recovery point objective for cross-region replication of the bucket. Applicable only for dual and multi-region buckets. `"DEFAULT"` sets default replication. `"ASYNC_TURBO"` value enables turbo replication, valid for dual-region buckets only. See [Turbo Replication](https://cloud.google.com/storage/docs/managing-turbo-replication) for more information. If rpo is not specified at bucket creation, it defaults to `"DEFAULT"` for dual and multi-region buckets. **NOTE** If used with single-region bucket, It will throw an error.

* `uniform_bucket_level_access` - (Optional, Default: false) Enables [Uniform bucket-level access](https://cloud.google.com/storage/docs/uniform-bucket-level-access) access to a bucket.

* `public_access_prevention` - (Optional) Prevents public access to a bucket. Acceptable values are "inherited" or "enforced". If "inherited", the bucket uses [public access prevention](https://cloud.google.com/storage/docs/public-access-prevention) only if the bucket is subject to the public access prevention organization policy constraint. Defaults to "inherited".

* `custom_placement_config` - (Optional) The bucket's custom location configuration, which specifies the individual regions that comprise a dual-region bucket. If the bucket is designated a single or multi-region, the parameters are empty. Structure is [documented below](#nested_custom_placement_config).

* `soft_delete_policy` -  (Optional, Computed) The bucket's soft delete policy, which defines the period of time that soft-deleted objects will be retained, and cannot be permanently deleted. If the block is not provided, Server side value will be kept which means removal of block won't generate any terraform change. Structure is [documented below](#nested_soft_delete_policy).

* `hierarchical_namespace` -  (Optional, ForceNew) The bucket's hierarchical namespace policy, which defines the bucket capability to handle folders in logical structure. Structure is [documented below](#nested_hierarchical_namespace). To use this configuration, `uniform_bucket_level_access` must be enabled on bucket.

* `time_created` -  (Computed) The creation time of the bucket in RFC 3339 format.

* `updated` -  (Computed) The time at which the bucket's metadata or IAM policy was last updated, in RFC 3339 format.

* `ip_filter` -  (Optional) The bucket IP filtering configuration. Specifies the network sources that can access the bucket, as well as its underlying objects. Structure is [documented below](#nested_ip_filter).

<a name="nested_lifecycle_rule"></a>The `lifecycle_rule` block supports:

* `action` - (Required) The Lifecycle Rule's action configuration. A single block of this type is supported. Structure is [documented below](#nested_action).

* `condition` - (Required) The Lifecycle Rule's condition configuration. A single block of this type is supported. Structure is [documented below](#nested_condition).

<a name="nested_action"></a>The `action` block supports:

* `type` - The type of the action of this Lifecycle Rule. Supported values include: `Delete`, `SetStorageClass` and `AbortIncompleteMultipartUpload`.

* `storage_class` - (Required if action type is `SetStorageClass`) The target [Storage Class](https://cloud.google.com/storage/docs/storage-classes) of objects affected by this Lifecycle Rule. Supported values include: `STANDARD`, `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`, `ARCHIVE`.

<a name="nested_condition"></a>The `condition` block supports the following elements, and requires at least one to be defined. If you specify multiple conditions in a rule, an object has to match all of the conditions for the action to be taken:

* `age` - (Optional) Minimum age of an object in days to satisfy this condition. **Note** To set `0` value of `age`, `send_age_if_zero` should be set `true` otherwise `0` value of `age` field will be ignored.

* `created_before` - (Optional) A date in the RFC 3339 format YYYY-MM-DD. This condition is satisfied when an object is created before midnight of the specified date in UTC.

* `with_state` - (Optional) Match to live and/or archived objects. Unversioned buckets have only live objects. Supported values include: `"LIVE"`, `"ARCHIVED"`, `"ANY"`.

* `matches_storage_class` - (Optional) [Storage Class](https://cloud.google.com/storage/docs/storage-classes) of objects to satisfy this condition. Supported values include: `STANDARD`, `MULTI_REGIONAL`, `REGIONAL`, `NEARLINE`, `COLDLINE`, `ARCHIVE`, `DURABLE_REDUCED_AVAILABILITY`.

* `matches_prefix` - (Optional) One or more matching name prefixes to satisfy this condition.

* `matches_suffix` - (Optional) One or more matching name suffixes to satisfy this condition.

* `num_newer_versions` - (Optional) Relevant only for versioned objects. The number of newer versions of an object to satisfy this condition. Due to a current bug you are unable to set this value to `0` within Terraform. When set to `0` it will be ignored and your state will treat it as though you supplied no `num_newer_versions` condition.

* `send_num_newer_versions_if_zero` - (Optional) While set true, `num_newer_versions` value will be sent in the request even for zero value of the field. This field is only useful for setting 0 value to the `num_newer_versions` field. It can be used alone or together with `num_newer_versions`.

* `custom_time_before` - (Optional) A date in the RFC 3339 format YYYY-MM-DD. This condition is satisfied when the customTime metadata for the object is set to an earlier date than the date used in this lifecycle condition.

* `days_since_custom_time` - (Optional)	Days since the date set in the `customTime` metadata for the object. This condition is satisfied when the current date and time is at least the specified number of days after the `customTime`. Due to a current bug you are unable to set this value to `0` within Terraform. When set to `0` it will be ignored, and your state will treat it as though you supplied no `days_since_custom_time` condition.

* `send_age_if_zero` - (Optional) While set true, `age` value will be sent in the request even for zero value of the field. This field is only useful and required for setting 0 value to the `age` field. It can be used alone or together with `age` attribute. **NOTE** `age` attibute with `0` value will be ommitted from the API request if `send_age_if_zero` field is having `false` value.

* `send_days_since_custom_time_if_zero` - (Optional) While set true, `days_since_custom_time` value will be sent in the request even for zero value of the field. This field is only useful for setting 0 value to the `days_since_custom_time` field. It can be used alone or together with `days_since_custom_time`.

* `days_since_noncurrent_time` - (Optional) Relevant only for versioned objects. Number of days elapsed since the noncurrent timestamp of an object. Due to a current bug you are unable to set this value to `0` within Terraform. When set to `0` it will be ignored, and your state will treat it as though you supplied no `days_since_noncurrent_time` condition.

* `send_days_since_noncurrent_time_if_zero` - (Optional) While set true, `days_since_noncurrent_time` value will be sent in the request even for zero value of the field. This field is only useful for setting 0 value to the `days_since_noncurrent_time` field. It can be used alone or together with `days_since_noncurrent_time`.

* `noncurrent_time_before` - (Optional) Relevant only for versioned objects. The date in RFC 3339 (e.g. `2017-06-13`) when the object became nonconcurrent. Due to a current bug you are unable to set this value to `0` within Terraform. When set to `0` it will be ignored, and your state will treat it as though you supplied no `noncurrent_time_before` condition.

<a name="nested_autoclass"></a>The `autoclass` block supports:

* `enabled` - (Required) While set to `true`, autoclass automatically transitions objects in your bucket to appropriate storage classes based on each object's access pattern.

* `terminal_storage_class` - (Optional) The storage class that objects in the bucket eventually transition to if they are not read for a certain length of time. Supported values include: `NEARLINE`, `ARCHIVE`.

<a name="nested_versioning"></a>The `versioning` block supports:

* `enabled` - (Required) While set to `true`, versioning is fully enabled for this bucket.

<a name="nested_website"></a>The `website` block supports the following elements, and requires at least one to be defined:

* `main_page_suffix` - (Optional) Behaves as the bucket's directory index where
    missing objects are treated as potential directories.

* `not_found_page` - (Optional) The custom object to return when a requested
    resource is not found.

<a name="nested_cors"></a>The `cors` block supports:

* `origin` - (Optional) The list of [Origins](https://tools.ietf.org/html/rfc6454) eligible to receive CORS response headers. Note: "*" is permitted in the list of origins, and means "any Origin".

* `method` - (Optional) The list of HTTP methods on which to include CORS response headers, (GET, OPTIONS, POST, etc) Note: "*" is permitted in the list of methods, and means "any method".

* `response_header` - (Optional) The list of HTTP headers other than the [simple response headers](https://www.w3.org/TR/cors/#simple-response-header) to give permission for the user-agent to share across domains.

* `max_age_seconds` - (Optional) The value, in seconds, to return in the [Access-Control-Max-Age header](https://www.w3.org/TR/cors/#access-control-max-age-response-header) used in preflight responses.

<a name="nested_retention_policy"></a>The `retention_policy` block supports:

* `is_locked` - (Optional) If set to `true`, the bucket will be [locked](https://cloud.google.com/storage/docs/using-bucket-lock#lock-bucket) and permanently restrict edits to the bucket's retention policy.  Caution: Locking a bucket is an irreversible action.

* `retention_period` - (Required) The period of time, in seconds, that objects in the bucket must be retained and cannot be deleted, overwritten, or archived. The value must be less than 2,147,483,647 seconds.

<a name="nested_logging"></a>The `logging` block supports:

* `log_bucket` - (Required) The bucket that will receive log objects.

* `log_object_prefix` - (Optional, Computed) The object prefix for log objects. If it's not provided,
    by default GCS sets this to this bucket's name.

<a name="nested_encryption"></a>The `encryption` block supports:

* `default_kms_key_name`: The `id` of a Cloud KMS key that will be used to encrypt objects inserted into this bucket, if no encryption method is specified.
  You must pay attention to whether the crypto key is available in the location that this bucket is created in.
  See [the docs](https://cloud.google.com/storage/docs/encryption/using-customer-managed-keys) for more details.

-> As per [the docs](https://cloud.google.com/storage/docs/encryption/using-customer-managed-keys) for customer-managed encryption keys, the IAM policy for the
  specified key must permit the [automatic Google Cloud Storage service account](https://cloud.google.com/storage/docs/projects#service-accounts) for the bucket's
  project to use the specified key for encryption and decryption operations.
  Although the service account email address follows a well-known format, the service account is created on-demand and may not necessarily exist for your project
  until a relevant action has occurred which triggers its creation.
  You should use the [`google_storage_project_service_account`](/docs/providers/google/d/storage_project_service_account.html) data source to obtain the email
  address for the service account when configuring IAM policy on the Cloud KMS key.
  This data source calls an API which creates the account if required, ensuring your Terraform applies cleanly and repeatedly irrespective of the
  state of the project.
  You should take care for race conditions when the same Terraform manages IAM policy on the Cloud KMS crypto key. See the data source page for more details.

<a name="nested_custom_placement_config"></a>The `custom_placement_config` block supports:

* `data_locations` - (Required) The list of individual regions that comprise a dual-region bucket. See [Cloud Storage bucket locations](https://cloud.google.com/storage/docs/dual-regions#availability) for a list of acceptable regions. **Note**: If any of the data_locations changes, it will [recreate the bucket](https://cloud.google.com/storage/docs/locations#key-concepts).

<a name="nested_soft_delete_policy"></a>The `soft_delete_policy` block supports:

* `retention_duration_seconds` - (Optional, Default: 604800) The duration in seconds that soft-deleted objects in the bucket will be retained and cannot be permanently deleted. Default value is 604800. The value must be in between 604800(7 days) and 7776000(90 days). **Note**: To disable the soft delete policy on a bucket, This field must be set to 0.

* `effective_time` - (Computed) Server-determined value that indicates the time from which the policy, or one with a greater retention, was effective. This value is in RFC 3339 format.

<a name="nested_hierarchical_namespace"></a>The `hierarchical_namespace` block supports:

* `enabled` - (Required) Enables hierarchical namespace for the bucket.

<a name="nested_ip_filter"></a>The `ip_filter` block supports:

* `mode` - (Required) The state of the IP filter configuration. Valid values are `Enabled` and `Disabled`. When set to `Enabled`, IP filtering rules are applied to a bucket and all incoming requests to the bucket are evaluated against these rules. When set to `Disabled`, IP filtering rules are not applied to a bucket. **Note**: `allow_all_service_agent_access` must be supplied when `mode` is set to `Enabled`, it can be ommited for other values.

* `allow_cross_org_vpcs` - (Optional) While set `true`, allows cross-org VPCs in the bucket's IP filter configuration.

* `allow_all_service_agent_access` (Optional) While set `true`, allows all service agents to access the bucket regardless of the IP filter configuration.

* `public_network_source` - (Optional) The public network IP address ranges that can access the bucket and its data. Structure is [documented below](#nested_public_network_source).

* `vpc_network_sources` - (Optional) The list of VPC networks that can access the bucket. Structure is [documented below](#nested_vpc_network_sources).

<a name="nested_public_network_source"></a>The `public_network_source` block supports:

* `allowed_ip_cidr_ranges` - The list of public IPv4 and IPv6 CIDR ranges that can access the bucket and its data.

<a name="nested_vpc_network_sources"></a>The `vpc_network_sources` block supports:

* `network` - Name of the network. Format: `projects/PROJECT_ID/global/networks/NETWORK_NAME`

* `allowed_ip_cidr_ranges` - The list of public or private IPv4 and IPv6 CIDR ranges that can access the bucket.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are
exported:

* `self_link` - The URI of the created resource.

* `url` - The base URL of the bucket, in the format `gs://<bucket-name>`.

## Timeouts

This resource provides the following
[Timeouts](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/retries-and-customizable-timeouts) configuration options: configuration options:

- `create` - Default is 4 minutes.
- `update` - Default is 4 minutes.
- `read` - Default is 4 minutes.

## Import

Storage buckets can be imported using the `name` or  `project/name`. If the project is not
passed to the import command it will be inferred from the provider block or environment variables.
If it cannot be inferred it will be queried from the Compute API (this will fail if the API is
not enabled).

* `{{project_id}}/{{bucket}}`
* `{{bucket}}`

In Terraform v1.5.0 and later, use an [`import` block](https://developer.hashicorp.com/terraform/language/import) to import Storage buckets using one of the formats above. For example:

```tf
import {
  id = "{{project_id}}/{{bucket}}"
  to = google_storage_bucket.default
}
```

When using the [`terraform import` command](https://developer.hashicorp.com/terraform/cli/commands/import), Storage buckets can be imported using one of the formats above. For example:

```
$ terraform import google_storage_bucket.default {{bucket}}
$ terraform import google_storage_bucket.default {{project_id}}/{{bucket}}
```

~> **Note:** Terraform will import this resource with `force_destroy` set to
`false` in state. If you've set it to `true` in config, run `terraform apply` to
update the value set in state. If you delete this resource before updating the
value, objects in the bucket will not be destroyed.
