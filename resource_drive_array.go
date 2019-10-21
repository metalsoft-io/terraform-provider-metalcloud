package metalcloud

import "log"
import "errors"
import "fmt"
import (
	"github.com/bigstepinc/metal-cloud-sdk-go"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceDriveArray() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"drive_array_label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"volume_template_label": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  nil,
			},
			"drive_array_storage_type": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Default:  "auto",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					if old == "auto" {
						return true
					}
					return false
				},
			},
			"drive_size_mbytes_default": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
				Default:  40960,
			},
		},
	}
}

func resourceDriveArrayCreate(infrastructureID int64, instanceArrayID int64, d map[string]interface{}, meta interface{}) error {

	client := meta.(*metalcloud.MetalCloudClient)

	var driveArrayStorageType = d["drive_array_storage_type"].(string)
	if driveArrayStorageType == "auto" {
		userLimits, err := client.InfrastructureUserLimits(infrastructureID)
		if err != nil {
			return err
		}
		driveArrayStorageTypesAvailable := (*userLimits)["storage_types"].([]interface{})
		driveArrayStorageType = driveArrayStorageTypesAvailable[0].(string)
	}

	//find the template id for the given volume_template_label and also build a list of possible values
	//in case we need to report it to the user
	availableVolumeTemplates, err := client.AvailableVolumeTemplatesGet()
	if err != nil {
		return err
	}

	var volumeTemplateID int64 = -1
	var possibleVolumeTemplateLabels []string

	for _, volumeTemplate := range *availableVolumeTemplates {
		if volumeTemplate.VolumeTemplateLabel == d["volume_template_label"].(string) {
			volumeTemplateID = volumeTemplate.VolumeTemplateID
			if volumeTemplate.VolumeTemplateDeprecationStatus != "not_deprecated" {
				log.Printf("WARNING: Template %s (%s) is DEPRECATED.", volumeTemplate.VolumeTemplateLabel, volumeTemplate.VolumeTemplateID)
			}
			break
		} else {
			possibleVolumeTemplateLabels = append(possibleVolumeTemplateLabels, volumeTemplate.VolumeTemplateLabel)
		}
	}

	if volumeTemplateID == -1 {
		return errors.New(
			fmt.Sprintf("Could not find template with volume_template_label=%s. Possible values are: %v.",
				d["volume_template_label"].(string),
				possibleVolumeTemplateLabels,
			))
	}

	driveArray := metalcloud.DriveArray{
		DriveArrayLabel:       d["drive_array_label"].(string),
		VolumeTemplateID:      volumeTemplateID,
		DriveArrayStorageType: driveArrayStorageType,
		InstanceArrayID:       instanceArrayID,
	}

	createdDriveArray, err2 := client.DriveArrayCreate(infrastructureID, driveArray)
	if err2 != nil {
		return err2
	}

	log.Printf("Created DriveArray %d", int(createdDriveArray.DriveArrayID))

	return nil
}

func resourceDriveArrayRead(driveArray metalcloud.DriveArray, meta interface{}) (*map[string]interface{}, error) {
	client := meta.(*metalcloud.MetalCloudClient)

	var ret = make(map[string]interface{})

	ret["drive_array_label"] = driveArray.DriveArrayLabel
	ret["drive_array_storage_type"] = driveArray.DriveArrayStorageType
	ret["drive_size_mbytes_default"] = int(driveArray.DriveSizeMBytesDefault)

	volumeTemplate, err := client.VolumeTemplateGet(driveArray.VolumeTemplateID)
	if err != nil {
		return nil, err
	}

	ret["volume_template_label"] = volumeTemplate.VolumeTemplateLabel

	return &ret, nil
}
