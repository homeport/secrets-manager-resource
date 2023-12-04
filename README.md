# Secrets Manager Resource

Concourse resource for secrets stored in IBM Cloud Secrets Manager instances.

## Source Configuration

- **endpointURL**: _Required_ Endpoint URL of the Secrets Manager instance to connect to, see [secrets manager docs](https://cloud.ibm.com/apidocs/secrets-manager/secrets-manager-v2?code=go#endpoints) for more details.
- **apikey**: _Required_ API key that allows access to read from the respective secrets manager instance.
- **secretName**: _Required_ Name of the secret in the secrets manager instance. This is the name, not the ID of the secret. The secret will be searched for by name through the API.

### Example

Since it is a custom resource type, it has to be configured once in the pipeline configuration.

```yaml
resource_types:
- name: secrets-manager-resource
  type: docker-image
  source:
   repository: ghcr.io/homeport/secrets-manager-resource
   tag: latest
```

One example would be to trigger a job, if the secret was updated in Secrets Manager.

```yaml
resources:
- name: some-secret
  type: secrets-manager-resource
  check_every: 2h
  icon: key
  source:
    endpointURL: https://<instance-id>.<region>.secrets-manager.appdomain.cloud
    apikey: ((your-api-key))
    secretName: super-important-secret

jobs:
- name: some-job
  plan:
  - get: some-secret
    trigger: true
    params:
      store-as: files
  - task: some-task
    config:
      inputs:
      - name: some-secret
      run:
        path: /bin/bash
        args:
        - -c
        - |
          #!/bin/bash
          some-tool login --secret $(< some-secret/payload)
```

## Behavior

### `check`: Checks for _updated at_ of a secret

Checks whether it finds a secret by the provided name and returns the last _updated at_ time.

### `in`: Obtains the secret data

Gets the secret by name and creates files based on the secret fields. Different secret types will create different files since they have different fields in Secret Manager. Check the [Working with secrets of different types](https://cloud.ibm.com/docs/secrets-manager?topic=secrets-manager-what-is-secret#secret-types) for more details on the types and their respective fields.

#### Parameters

* `store-as`: *Optional.* Defines what to do with the secret.
  * `file` (or `files`) will use the API to obtain the secrets entry and store the content on disk with one file per secret entry detail, e.g. files like `name`, `description`, or `payload`. Please note, fields like `payload` will contain Base64 encoded content. Use this option with caution since the secret data will be stored on disk for the runtime of the container.
  * `script` will create a `get-secret.sh` shell script that contains the `curl` command to obtain the secrets details. You have to make sure to be logged in the right account before calling the generated script so that the correct access tokens can be requested.
  * _empty_/_omitted_ will result in no look-up of the secret itself and no data is written. Use this if you only need to be notified about the update of the secret details.

### `out`: No-op

Not implemented. May be subject to change in the future.

## Development

### Prerequisites

- Go is _Required_ - version 1.20 is in use, newer versions will probably work.
- Docker or similar is _Required_ - any tool that allows for a `docker build` like container build.

### Contributing

Please make all pull requests to the `main` branch and ensure tests pass locally.
