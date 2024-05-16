package metalcloud

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
)

// DataSourceInfrastructure provides a way to search among existing infrastructures and create if not exists
func DataSourceInfrastructureReference() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInfrastructureReferenceRead,
		Schema: map[string]*schema.Schema{
			"infrastructure_label": {
				Type:     schema.TypeString,
				Required: true,
				//this required as the serverside will convert to lowercase and generate a diff
				//also helpful to prevent other
				ValidateDiagFunc: validateLabel,
				ForceNew:         true,
				DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
					if strings.ToLower(old) == strings.ToLower(new) {
						return true
					}

					if new == "" {
						return true
					}
					return false
				},
			},
			"datacenter_name": {
				Type:     schema.TypeString,
				Required: true,
				//this required as the serverside will convert to lowercase and generate a diff
				//also helpful to prevent other
				ValidateDiagFunc: validateLabel,
			},
			"create_if_not_exists": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"infrastructure_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceInfrastructureReferenceRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	var diags diag.Diagnostics

	client := meta.(*mc.Client)

	//we list all available infrastructures to the currently logged in user
	//if the one we are looking forward exists we return it
	//we do this instead of infrastructureGetByLabel because we don't have good error codes to distinguish between
	//a not-found error and another type of error
	infrastructures, err := client.Infrastructures()
	if err != nil {
		return diag.FromErr(err)
	}

	infrastructure_label := d.Get("infrastructure_label").(string)
	datacenter_name := d.Get("datacenter_name").(string)

	//if the one we are looking forward exists we return it
	iRet, ok := (*infrastructures)[strings.ToLower(infrastructure_label)]
	if ok {
		//assert if datacenter name matches
		if datacenter_name != iRet.DatacenterName {
			return diag.Errorf("Datacenter of infrastructure '%s' returned from the server '%s' is different from the one defined on the datasource'%s'", infrastructure_label, iRet.DatacenterName, datacenter_name)
		}
	} else if d.Get("create_if_not_exists").(bool) {
		//if could not find it we create it
		i := mc.Infrastructure{
			InfrastructureLabel: infrastructure_label,
			DatacenterName:      datacenter_name,
		}

		obj, err := client.InfrastructureCreate(i)
		if err != nil {
			return diag.FromErr(err)
		}

		iRet = *obj
	} else {
		return diag.Errorf("Infrastructure \"%s\" does not exist and you specified not to create it.", infrastructure_label)
	}

	d.SetId(fmt.Sprintf("%d", iRet.InfrastructureID))
	d.Set("infrastructure_id", iRet.InfrastructureID)

	return diags

}
