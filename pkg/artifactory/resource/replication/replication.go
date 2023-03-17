package replication

import (
	"context"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jfrog/terraform-provider-shared/client"
)

const EndpointPath = "artifactory/api/replications/"

var replicationSchemaEnableEventReplication = map[string]*schema.Schema{
	"enable_event_replication": {
		Type:        schema.TypeBool,
		Optional:    true,
		Default:     false,
		Description: "When set, each event will trigger replication of the artifacts changed in this event. This can be any type of event on artifact, e.g. add, deleted or property change. Default value is `false`.",
	},
}

func resourceReplicationDelete(_ context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	_, err := m.(*resty.Client).R().
		AddRetryCondition(client.RetryOnMergeError).
		Delete(EndpointPath + d.Id())
	return diag.FromErr(err)
}

type repoConfiguration struct {
	Rclass string `json:"rclass"`
}

func getRepositoryRclass(repoKey string, m interface{}) (string, error) {
	repoConfig := repoConfiguration{}
	_, err := m.(*resty.Client).R().
		SetResult(&repoConfig).
		Get("artifactory/api/repositories/" + repoKey)
	if err != nil {
		return "", err
	}
	return repoConfig.Rclass, err
}

func verifyRepoRclass(repoKey string, expectedRclass string, m interface{}) (bool, error) {
	rclass, err := getRepositoryRclass(repoKey, m)
	if err != nil {
		return false, fmt.Errorf("error getting repository configuration: %v", err)
	}
	if rclass == expectedRclass {
		return true, nil
	}
	return false, nil
}
