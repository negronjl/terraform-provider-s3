$TERRAFORM_PLUGINS_DIR="$($env:APPDATA)\terraform.d\plugins\windows_amd64" 
$PROVIDER_PATH=(Join-Path ${TERRAFORM_PLUGINS_DIR} "terraform-provider-s3.exe")
go build
go install
Copy-Item "$($env:GOPATH)\bin\terraform-provider-s3.exe" $PROVIDER_PATH
