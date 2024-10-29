package metalcloud

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	mc "github.com/metalsoft-io/metal-cloud-sdk-go/v3"
	mc2 "github.com/metalsoft-io/metal-cloud-sdk2-go"
)

func resourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServerCreate,
		ReadContext:   resourceServerRead,
		UpdateContext: resourceServerUpdate,
		DeleteContext: resourceServerDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"server_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"server_type_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"datacenter_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_uuid": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"serial_number": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vendor": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"vendor_sku_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"model": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"submodel": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bmc_mac_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bmc_hostname": {
				Type:     schema.TypeString,
				Required: true,
			},
			"bmc_username": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"bmc_password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"impi_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"processor_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"processor_core_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"processor_core_mhz": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"processor_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"processor_cpu_mark": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"processor_threads": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"gpu_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"ram_gbytes": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"disk_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"server_disk_count": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"capacity_mbps": {
				Type:     schema.TypeFloat,
				Optional: true,
				Computed: true,
			},
			"mgmt_snmp_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"vnc_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"supports_sol": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"supports_virtual_media": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"supports_oob_provisioning": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"supports_fc_provisioning": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"is_endpoint": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"cleanup_policy_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"server_comments": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"requires_re_register": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"server_disk_wipe": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"administration_state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"server_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"dhcp_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"power_status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client2, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	client := meta.(*mc.Client)

	server := expandServerRegistration(d, client)

	newServer, _, err := client2.ServerApi.RegisterServer(ctx, server)
	if err != nil {
		return extractApiError(err)
	}

	id := fmt.Sprintf("%d", newServer.Id)
	d.SetId(id)

	return resourceServerRead(ctx, d, meta)
}

func resourceServerRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client2, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	server, _, err := client2.ServerApi.GetServerInfo(ctx, float64(id))
	if err != nil {
		return extractApiError(err)
	}

	flattenServer(d, server)

	return diags
}

func resourceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client2, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	oldServer, _, err := client2.ServerApi.GetServerInfo(ctx, float64(id))
	if err != nil {
		return extractApiError(err)
	}

	newServer := expandServer(d)

	if newServer.PowerStatus != oldServer.PowerStatus {
		powerSet := mc2.ServerPowerSetDto{
			PowerCommand: newServer.PowerStatus,
		}

		_, err = client2.ServerApi.SetServerPowerState(ctx, powerSet, float64(id))
		if err != nil {
			return extractApiError(err)
		}
	}

	return resourceServerRead(ctx, d, meta)
}

func resourceServerDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	client2, err := getClient2()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, _, err = client2.ServerApi.GetServerInfo(ctx, float64(id))

	// if the server has already been deleted we ignore it, else we actually delete because that might have been the intent.
	if err == nil {
		client := meta.(*mc.Client)

		err = client.SecretDelete(id)
		if err != nil {
			return extractApiError(err)
		}
	}

	d.SetId("")
	return diags
}

func flattenServer(d *schema.ResourceData, server mc2.Server) map[string]interface{} {
	d.Set("server_id", server.ServerId)
	d.Set("server_type_id", server.ServerTypeId)
	d.Set("datacenter_name", server.DatacenterName)
	d.Set("server_uuid", server.ServerUUID)
	d.Set("serial_number", server.SerialNumber)
	d.Set("bmc_mac_address", server.BmcMacAddress)
	d.Set("vendor", server.Vendor)
	d.Set("vendor_sku_id", server.VendorSkuId)
	d.Set("model", server.Model)
	d.Set("submodel", server.Submodel)
	d.Set("bmc_hostname", server.ManagementAddress)
	d.Set("bmc_username", server.Username)
	d.Set("impi_version", server.ImpiVersion)
	d.Set("processor_count", server.ProcessorCount)
	d.Set("processor_core_count", server.ProcessorCoreCount)
	d.Set("processor_core_mhz", server.ProcessorCoreMhz)
	d.Set("processor_name", server.ProcessorName)
	d.Set("processor_cpu_mark", server.ProcessorCpuMark)
	d.Set("processor_threads", server.ProcessorThreads)
	d.Set("gpu_count", server.GpuCount)
	d.Set("ram_gbytes", server.RamGbytes)
	d.Set("disk_count", server.DiskCount)
	d.Set("server_disk_count", server.ServerDiskCount)
	d.Set("capacity_mbps", server.ServerCapacityMbps)
	d.Set("mgmt_snmp_port", server.MgmtSnmpPort)
	d.Set("vnc_port", server.VncPort)
	d.Set("supports_sol", server.ServerSupportsSol)
	d.Set("supports_virtual_media", server.ServerSupportsVirtualMedia)
	d.Set("supports_oob_provisioning", server.ServerSupportsOobProvisioning)
	d.Set("supports_fc_provisioning", server.SupportsFcProvisioning)
	d.Set("is_endpoint", server.IsBasicCampusEndpoint)
	d.Set("cleanup_policy_id", server.ServerCleanupPolicyId)
	d.Set("server_comments", server.ServerComments)
	d.Set("requires_re_register", server.RequiresReRegister)
	d.Set("server_disk_wipe", server.ServerDiskWipe)
	d.Set("administration_state", server.AdministrationState)
	d.Set("server_status", server.ServerStatus)
	d.Set("dhcp_status", server.ServerDhcpStatus)
	d.Set("power_status", server.PowerStatus)

	return nil
}

func expandServer(d *schema.ResourceData) mc2.Server {
	var server mc2.Server

	if d.Get("server_id") != nil {
		server.ServerId = d.Get("server_id").(float64)
	}

	server.ServerTypeId = d.Get("server_type_id").(float64)
	server.DatacenterName = d.Get("datacenter_name").(string)
	server.ServerUUID = d.Get("server_uuid").(string)
	server.SerialNumber = d.Get("serial_number").(string)
	server.BmcMacAddress = d.Get("bmc_mac_address").(string)
	server.Vendor = d.Get("vendor").(string)
	server.VendorSkuId = d.Get("vendor_sku_id").(string)
	server.Model = d.Get("model").(string)
	server.Submodel = d.Get("submodel").(string)
	server.ManagementAddress = d.Get("bmc_hostname").(string)
	server.Username = d.Get("bmc_username").(string)
	server.ImpiVersion = d.Get("impi_version").(string)
	server.ProcessorCount = d.Get("processor_count").(float64)
	server.ProcessorCoreCount = d.Get("processor_core_count").(float64)
	server.ProcessorCoreMhz = d.Get("processor_core_mhz").(float64)
	server.ProcessorName = d.Get("processor_name").(string)
	server.ProcessorCpuMark = d.Get("processor_cpu_mark").(float64)
	server.ProcessorThreads = d.Get("processor_threads").(float64)
	server.GpuCount = d.Get("gpu_count").(float64)
	server.RamGbytes = d.Get("ram_gbytes").(float64)
	server.DiskCount = d.Get("disk_count").(float64)
	server.ServerDiskCount = d.Get("server_disk_count").(float64)
	server.ServerCapacityMbps = d.Get("capacity_mbps").(float64)
	server.MgmtSnmpPort = d.Get("mgmt_snmp_port").(float64)
	server.VncPort = d.Get("vnc_port").(float64)
	server.ServerSupportsSol = d.Get("supports_sol").(float64)
	server.ServerSupportsVirtualMedia = d.Get("supports_virtual_media").(float64)
	server.ServerSupportsOobProvisioning = d.Get("supports_oob_provisioning").(float64)
	server.SupportsFcProvisioning = d.Get("supports_fc_provisioning").(float64)
	server.IsBasicCampusEndpoint = d.Get("is_endpoint").(float64)
	server.ServerCleanupPolicyId = d.Get("cleanup_policy_id").(float64)
	server.RequiresReRegister = d.Get("requires_re_register").(float64)
	server.ServerComments = d.Get("server_comments").(string)
	server.ServerDiskWipe = d.Get("server_disk_wipe").(float64)
	server.AdministrationState = d.Get("administration_state").(string)
	server.ServerStatus = d.Get("server_status").(string)
	server.ServerDhcpStatus = d.Get("dhcp_status").(string)
	server.PowerStatus = d.Get("power_status").(string)

	return server
}

func expandServerRegistration(d *schema.ResourceData, client *mc.Client) mc2.ServerRegistrationDto {
	var server mc2.ServerRegistrationDto

	server.BmcHostname = d.Get("bmc_hostname").(string)
	server.BmcUser = d.Get("bmc_username").(string)
	server.BmcPassword = d.Get("bmc_password").(string)
	server.Vendor = d.Get("vendor").(string)
	server.Model = d.Get("model").(string)
	server.MacAddress = d.Get("bmc_mac_address").(string)
	server.SerialNumber = d.Get("serial_number").(string)
	server.Uuid = d.Get("server_uuid").(string)

	datacenterName := d.Get("datacenter_name").(string)
	if datacenterName != "" {
		datacenter, err := client.DatacenterGet(datacenterName)
		if err == nil {
			server.SiteId = float64(datacenter.DatacenterID)
		}
	}

	return server
}
