package metalcloud

import (
	"fmt"
	"log"

	mc "github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

//DataSourceVolumeTemplate provides means to search for volume templates
func DataSourceVolumeTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceVolumeTemplateRead,
		Schema: map[string]*schema.Schema{
			"volume_template_label": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"volume_template_id": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceVolumeTemplateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*mc.Client)

	volumeTemplateLabel := d.Get("volume_template_label").(string)
	//find the template id for the given volume_template_label and also build a list of possible values
	//in case we need to report it to the user
	availableVolumeTemplates, err := client.VolumeTemplates()
	if err != nil {
		return err
	}

	var volumeTemplateID = -1
	var possibleVolumeTemplateLabels []string

	for _, volumeTemplate := range *availableVolumeTemplates {
		if volumeTemplate.VolumeTemplateLabel == volumeTemplateLabel {
			volumeTemplateID = volumeTemplate.VolumeTemplateID
			if volumeTemplate.VolumeTemplateDeprecationStatus != "not_deprecated" {
				log.Printf("WARNING: Template %s (%d) is DEPRECATED.", volumeTemplate.VolumeTemplateLabel, volumeTemplate.VolumeTemplateID)
			}
			break
		} else {
			possibleVolumeTemplateLabels = append(possibleVolumeTemplateLabels, volumeTemplate.VolumeTemplateLabel)
		}
	}

	if volumeTemplateID == -1 {
		return fmt.Errorf("Could not find template with volume_template_label=%s. Possible values are: %v",
			volumeTemplateLabel,
			possibleVolumeTemplateLabels)
	}

	log.Printf("VolumeTemplateID = %d", volumeTemplateID)
	d.Set("volume_template_id", volumeTemplateID)
	d.SetId(fmt.Sprintf("%d", volumeTemplateID))

	return nil
}
