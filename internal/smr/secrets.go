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

	sm "github.com/IBM/secrets-manager-go-sdk/v2/secretsmanagerv2"
	"github.com/go-openapi/strfmt"
)

// SecretMetadata is slimmed down version that _should_ match all of the possible
// secret types as it only contains the very common fields.
type SecretMetadata struct {
	CreatedBy        *string                `json:"created_by,omitempty"`
	CreatedAt        *strfmt.DateTime       `json:"created_at,omitempty"`
	Crn              *string                `json:"crn,omitempty"`
	CustomMetadata   map[string]interface{} `json:"custom_metadata,omitempty"`
	Description      *string                `json:"description,omitempty"`
	Downloaded       *bool                  `json:"downloaded,omitempty"`
	ID               *string                `json:"id,omitempty"`
	Labels           []string               `json:"labels,omitempty"`
	LocksTotal       *int64                 `json:"locks_total,omitempty"`
	Name             *string                `json:"name,omitempty"`
	SecretGroupID    *string                `json:"secret_group_id,omitempty"`
	SecretType       *string                `json:"secret_type,omitempty"`
	State            *int64                 `json:"state,omitempty"`
	StateDescription *string                `json:"state_description,omitempty"`
	UpdatedAt        *strfmt.DateTime       `json:"updated_at,omitempty"`
}

func (s *SecretMetadata) Id() (string, error) {
	if s.ID == nil {
		return "", fmt.Errorf("no ID defined")
	}

	return *s.ID, nil
}

func GetSecretMetaDataBySecretName(service sm.SecretsManagerV2, source Source) (*SecretMetadata, error) {
	listSecretsOptions := &sm.ListSecretsOptions{Search: &source.SecretName}

	if source.SecretGroupID != "" {
		listSecretsOptions.Groups = append(listSecretsOptions.Groups, source.SecretGroupID)
	}

	pager, err := service.NewSecretsPager(listSecretsOptions)
	if err != nil {
		return nil, err
	}

	var results []sm.SecretMetadataIntf
	for pager.HasNext() {
		nextPage, err := pager.GetNext()
		if err != nil {
			return nil, err
		}

		results = append(results, nextPage...)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("cannot find secret with name %q", source.SecretName)
	}

	if len(results) != 1 {
		return nil, fmt.Errorf("more than one secret was found searching for %q", source.SecretName)
	}

	data, err := json.Marshal(results[0])
	if err != nil {
		return nil, err
	}

	var result SecretMetadata
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
