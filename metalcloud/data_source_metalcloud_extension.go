package metalcloud

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	log "github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	sdk2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
)

func DataSourceExtension() *schema.Resource {
	return &schema.Resource{
		SchemaFunc:  dataSourceExtensionSchema,
		ReadContext: dataSourceExtensionRead,
	}
}

func dataSourceExtensionSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		fieldExtensionId: {
			Type:     schema.TypeInt,
			Optional: true,
			Computed: true,
		},
		fieldExtensionLabel: {
			Type:     schema.TypeString,
			Required: true,
			DiffSuppressFunc: func(_, old, new string, d *schema.ResourceData) bool {
				return strings.EqualFold(old, new)
			},
		},
	}
}

func dataSourceExtensionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	extension_id := d.Get(fieldExtensionId).(int)
	if extension_id != 0 && !d.HasChange(fieldExtensionLabel) {
		x, err := getExtension(ctx, extension_id)
		if err != nil {
			if errors.Is(err, errNotFound) {
				d.SetId("")
			}

			return extractApiError(err)
		}

		d := setFromExtensionData(d, int(x.Id), x.Label, x.Status)
		if d != nil {
			diags = append(diags, *d)
		}
	} else {
		label := d.Get(fieldExtensionLabel).(string)

		x, err := findExtension(ctx, label)
		if err != nil {
			return extractApiError(err)
		}

		if x != nil {
			d := setFromExtensionData(d, int(x.Id), x.Label, x.Status)
			if d != nil {
				diags = append(diags, *d)
			}
		} else {
			d.SetId("")
		}
	}

	return diags
}

func getExtension(ctx context.Context, id int) (*sdk2.ExtensionDto, error) {
	client, err := getClient2()
	if err != nil {
		return nil, err
	}

	x, r, err :=
		client.ExtensionApi.GetExtension(ctx, float64(id))
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v: %v\r\n", r.StatusCode, err))

		if r.StatusCode == http.StatusNotFound {
			return nil, errNotFound
		}

		return nil, err
	}

	return &x, err
}

func findExtension(ctx context.Context, label string) (*sdk2.ExtensionInfoDto, error) {
	client, err := getClient2()
	if err != nil {
		return nil, err
	}

	xTypes, r, err := client.ExtensionApi.GetExtensions(ctx, nil)
	if err != nil {
		log.Debug(ctx, fmt.Sprintf("%v: %v\r\n", r.StatusCode, err))

		return nil, err
	}

	for _, x := range xTypes.Extensions {
		if strings.EqualFold(x.Label, label) {
			return &x, nil
		}
	}

	return nil, err
}

func setFromExtensionData(d *schema.ResourceData, id int, label string, status string) *diag.Diagnostic {
	d.Set(fieldExtensionId, id)
	d.Set(fieldExtensionLabel, label)

	d.SetId(fmt.Sprintf("%d", id))

	if status != extensionStatus_Active {
		return &diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  "Not active",
			Detail:   fmt.Sprintf("Extension cannot be used for instantiation of new instances. Status: %v.", status),
		}
	}

	return nil
}
