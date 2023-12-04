// Copyright © 2023 The Homeport Team
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package smr

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/gonvenience/bunt"

	sm "github.com/IBM/secrets-manager-go-sdk/secretsmanagerv2"
)

const getSecretScript = `#!/bin/sh

set -e
set -u

curl \
  --fail \
  --silent \
  --request GET \
  --location \
  --header "Authorization: $(ibmcloud iam oauth-tokens --output json | jq --raw-output .iam_token)" \
  --header "Accept: application/json" \
  "%s/api/v2/secrets/%s"

`

func Check(r io.Reader) error {
	config, err := LoadConfig(r)
	if err != nil {
		return err
	}

	secretsManagerService, err := sm.NewSecretsManagerV2(config.SecretsManagerV2Options())
	if err != nil {
		return err
	}

	metadata, err := GetSecretMetaDataBySecretName(*secretsManagerService, config.Source.SecretName)
	if err != nil {
		return err
	}

	bunt.Fprintf(os.Stderr, "Found secret with name Yellow{%s} of type Orange{%s}, created at PaleGreen{%v}, updated at ForestGreen{%v}.\n",
		*metadata.Name,
		*metadata.SecretType,
		metadata.CreatedAt,
		metadata.UpdatedAt,
	)

	out, err := json.Marshal([]Version{{Timestamp: time.Time(*metadata.UpdatedAt)}})
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stdout, string(out))
	return err
}

func In(r io.Reader, target string) error {
	config, err := LoadInConfig(r)
	if err != nil {
		return err
	}

	secretsManagerService, err := sm.NewSecretsManagerV2(config.SecretsManagerV2Options())
	if err != nil {
		return err
	}

	metadata, err := GetSecretMetaDataBySecretName(*secretsManagerService, config.Source.SecretName)
	if err != nil {
		return err
	}

	bunt.Fprintf(os.Stderr, "Found secret with name Yellow{%s} of type Orange{%s}, created at PaleGreen{%v}, updated at ForestGreen{%v}.\n",
		*metadata.Name,
		*metadata.SecretType,
		metadata.CreatedAt,
		metadata.UpdatedAt,
	)

	if !time.Time(*metadata.UpdatedAt).Equal(config.Version.Timestamp) {
		return fmt.Errorf("mismatch between config timestamp (%v) and secret timestamp (%v)",
			config.Version.Timestamp,
			metadata.UpdatedAt,
		)
	}

	id, err := metadata.Id()
	if err != nil {
		return err
	}

	switch config.Params["store-as"] {
	case "file", "files":
		secret, _, err := secretsManagerService.GetSecret(secretsManagerService.NewGetSecretOptions(id))
		if err != nil {
			return err
		}

		b, err := json.Marshal(secret)
		if err != nil {
			return err
		}

		var mapping map[string]interface{}
		if err := json.Unmarshal(b, &mapping); err != nil {
			return err
		}

		for k, v := range mapping {
			switch obj := v.(type) {
			case string:
				if err := os.WriteFile(filepath.Join(target, k), []byte(obj), 0644); err != nil {
					return err
				}
			}
		}

	case "script":
		var script = fmt.Sprintf(getSecretScript, config.Source.EndpointURL, id)
		if err := os.WriteFile(filepath.Join(target, "get-secret.sh"), []byte(script), 0755); err != nil {
			return err
		}

	default:
		bunt.Fprintf(os.Stderr, "no or unknown _store-as_ parameter defined, no files are created\n")
	}

	// --- --- ---

	output := Output{
		Version: Version{Timestamp: time.Time(*metadata.UpdatedAt)},
		Metadata: []KeyVal{
			{"Name", *metadata.Name},
			{"SecretType", *metadata.SecretType},
			{"StateDescription", *metadata.StateDescription},
			{"CreatedAt", metadata.CreatedAt.String()},
			{"UpdatedAt", metadata.UpdatedAt.String()},
			{"SecretId", id},
		},
	}

	out, err := json.Marshal(output)
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(os.Stdout, string(out))
	return err
}
