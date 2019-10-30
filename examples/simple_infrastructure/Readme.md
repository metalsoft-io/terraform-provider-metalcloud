##simple example of metalcloud terraform

Set key & api variables:

```bash
export TF_VAR_api_key="<yourkey>"
export TF_VAR_user_email="test@test.com"
export TF_VAR_endpoint="https://api.bigstep.com/metal-cloud"
export TF_VAR_datacenter="uk-reading"
```

The plan phase:
```bash
terraform plan
```

The apply phase:
```bash
terraform apply
```

To delete the infrastrucure:
```bash
terraform destroy
```


