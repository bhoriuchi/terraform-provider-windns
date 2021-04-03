package windns

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"

	"github.com/masterzen/winrm"
)

// ReadARecordOptions options to read an a record
type ReadARecordOptions struct {
	DnsServer string
	Name      string
	Address   string
	ZoneName  string
}

// AddARecordOptions options to add an a record
type AddARecordOptions struct {
	DnsServer      string
	Name           string
	Address        string
	ZoneName       string
	AllowUpdateAny bool
	CreatePtr      bool
	TTL            int
}

// UpdateARecordOptions options to add an a record
type UpdateARecordOptions struct {
	DnsServer  string
	Name       string
	Address    string
	NewAddress string
	ZoneName   string
	TTL        int
}

// DeleteARecordOptions options to add an a record
type DeleteARecordOptions struct {
	DnsServer string
	Name      string
	Address   string
	ZoneName  string
}

// ReadARecord reads an A record
func (c *Client) ReadARecord(opts *ReadARecordOptions) (*Response, error) {
	if opts.Name == "" {
		return nil, fmt.Errorf(`required value "name" not spcified`)
	}
	if opts.ZoneName == "" {
		return nil, fmt.Errorf(`required value "zone_name" not spcified`)
	}

	opts.DnsServer = c.o.DnsServer
	w := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	tpl := template.Must(template.New("main").Parse(readARecordScript))
	if err := tpl.Execute(w, opts); err != nil {
		return nil, err
	}

	ps := winrm.Powershell(w.String())
	exitCode, err := c.c.Run(ps, stdout, stderr)
	if err != nil {
		return nil, err
	} else if exitCode != 0 {
		return nil, fmt.Errorf("exit code %d: %s", exitCode, stderr.String())
	}

	rsp := &Response{}
	if err := json.Unmarshal(stdout.Bytes(), rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

// AddARecord adds a new A record
func (c *Client) AddARecord(opts *AddARecordOptions) (*Response, error) {
	if opts.Name == "" {
		return nil, fmt.Errorf(`required value "name" not spcified`)
	}
	if opts.Address == "" {
		return nil, fmt.Errorf(`required value "address" not spcified`)
	}
	if opts.ZoneName == "" {
		return nil, fmt.Errorf(`required value "zone_name" not spcified`)
	}

	if opts.TTL < 1 {
		opts.TTL = defaultTTL
	}

	opts.DnsServer = c.o.DnsServer
	w := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	tpl := template.Must(template.New("main").Parse(addARecordScript))
	if err := tpl.Execute(w, opts); err != nil {
		return nil, err
	}

	ps := winrm.Powershell(w.String())
	exitCode, err := c.c.Run(ps, stdout, stderr)
	if err != nil {
		return nil, err
	} else if exitCode != 0 {
		return nil, fmt.Errorf("exit code %d: %s", exitCode, stderr.String())
	}

	rsp := &Response{}
	if err := json.Unmarshal(stdout.Bytes(), rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

// UpdateARecord adds an A record
func (c *Client) UpdateARecord(opts *UpdateARecordOptions) (*Response, error) {
	if opts.Name == "" {
		return nil, fmt.Errorf(`required value "name" not spcified`)
	}
	if opts.ZoneName == "" {
		return nil, fmt.Errorf(`required value "zone_name" not spcified`)
	}
	if opts.Address == "" {
		return nil, fmt.Errorf(`required value "address" not spcified`)
	}

	opts.DnsServer = c.o.DnsServer
	w := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	tpl := template.Must(template.New("main").Parse(updateARecordScript))
	if err := tpl.Execute(w, opts); err != nil {
		return nil, err
	}

	ps := winrm.Powershell(w.String())
	exitCode, err := c.c.Run(ps, stdout, stderr)
	if err != nil {
		return nil, err
	} else if exitCode != 0 {
		return nil, fmt.Errorf("exit code %d: %s", exitCode, stderr.String())
	}

	rsp := &Response{}
	if err := json.Unmarshal(stdout.Bytes(), rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

// DeleteARecord deletes an A record
func (c *Client) DeleteARecord(opts *DeleteARecordOptions) (*Response, error) {
	if opts.Name == "" {
		return nil, fmt.Errorf(`required value "name" not spcified`)
	}
	if opts.ZoneName == "" {
		return nil, fmt.Errorf(`required value "zone_name" not spcified`)
	}
	if opts.Address == "" {
		return nil, fmt.Errorf(`required value "address" not spcified`)
	}

	opts.DnsServer = c.o.DnsServer
	w := new(bytes.Buffer)
	stdout := new(bytes.Buffer)
	stderr := new(bytes.Buffer)
	tpl := template.Must(template.New("main").Parse(deleteARecordScript))
	if err := tpl.Execute(w, opts); err != nil {
		return nil, err
	}

	ps := winrm.Powershell(w.String())
	exitCode, err := c.c.Run(ps, stdout, stderr)
	if err != nil {
		return nil, err
	} else if exitCode != 0 {
		return nil, fmt.Errorf("exit code %d: %s", exitCode, stderr.String())
	}

	rsp := &Response{}
	if err := json.Unmarshal(stdout.Bytes(), rsp); err != nil {
		return nil, err
	}

	return rsp, nil
}

const (
	readARecordScript = `
	Import-Module DNSServer

	$findArgs = @{
		Name         = "{{.Name}}"
		ComputerName = "{{.DnsServer}}"
		ZoneName     = "{{.ZoneName}}"
		RRType       = "A"
		ErrorAction  = "SilentlyContinue"
	}
		
	$record = Get-DnsServerResourceRecord @findArgs
	if ($Error.Count -gt 0) {
		if ($Error[0].CategoryInfo.Category -eq "ObjectNotFound")
		{
			$res = @{
							code = 404
							detail = "record not found"
					}
			Write-Output "$($res | ConvertTo-Json -Compress)"
			return
		}
		else {
				$res = @{
							code = 500
							detail = "$($Error[0].Exception.Message)"
					}
				Write-Output "$($res | ConvertTo-Json -Compress)"
				return
		}
	}

	$records = @()
	$record | ForEach-Object {
		$addr = $_.RecordData.IPv4Address.IPAddressToString
		if (-not (![string]::IsNullOrEmpty("{{.Address}}") -and "{{.Address}}" -ne $addr)) {
			$records += @{
				type = "A"
				name = $_.HostName
				data = $addr
				zone = "{{.ZoneName}}"
				ttl  = $_.TimeToLive.TotalSeconds
			}
		}
	}
	
	$res = @{
			code = 200
			detail  = "record found"
			records = $records
	}

	Write-Output "$($res | ConvertTo-Json -Compress -Depth 5)"
	`

	addARecordScript = `
	Import-Module DNSServer

	$findArgs = @{
		Name         = "{{.Name}}"
		ComputerName = "{{.DnsServer}}"
		ZoneName     = "{{.ZoneName}}"
		RRType       = "A"
		ErrorAction  = "SilentlyContinue"
	}
	
	$record = Get-DnsServerResourceRecord @findArgs
	if ($Error.Count -gt 0) {
		if ($Error[0].CategoryInfo.Category -eq "ObjectNotFound")
		{
			$Error.Clear()
		}
		else {
				$res = @{
							code = 500
							detail = "$($Error[0].Exception.Message)"
					}
				Write-Output "$($res | ConvertTo-Json -Compress)"
				return
		}
	}

	if (![string]::IsNullOrEmpty("{{.Address}}")) {
		$record = $record | Where-Object {
			$_.RecordData.IPv4Address.IPAddressToString -eq "{{.Address}}"
		}
	}
	
	if ($null -ne $record) {
		$res = @{
					code = 400
					detail = "record already exists"
			}
		Write-Output "$($res | ConvertTo-Json -Compress)"
		return
	}
	
	$createArgs = @{
		A              = $true
		ZoneName       = "{{.ZoneName}}"
		Name           = "{{.Name}}"
		IPv4Address    = "{{.Address}}"
		ComputerName   = "{{.DnsServer}}"
		AllowUpdateAny = ${{.AllowUpdateAny}}
		CreatePtr      = ${{.CreatePtr}}
		TimeToLive     = [System.TimeSpan]::FromSeconds({{.TTL}})
		Confirm        = $false
		ErrorAction    = "SilentlyContinue"
	}
	
	Add-DnsServerResourceRecord @createArgs
	if ($Errors.Count -gt 0) {
		$res = @{
					code = 500
					detail = "$($Error[0].Exception.Message)"
			}
		Write-Output "$($res | ConvertTo-Json -Compress)"
		return
	}
	
	$records = @()
	$records += @{
		type = "A"
		name = "{{.Name}}"
		data = "{{.Address}}"
		zone = "{{.ZoneName}}"
		ttl  = {{.TTL}}
	}

	$res = @{
			code = 200
			detail = "record created"
			records = $records
	}
	
	Write-Output "$($res | ConvertTo-Json -Compress -Depth 5)"
	`

	updateARecordScript = `
	Import-Module DNSServer

	$findArgs = @{
		Name         = "{{.Name}}"
		ComputerName = "{{.DnsServer}}"
		ZoneName     = "{{.ZoneName}}"
		RRType       = "A"
		ErrorAction  = "SilentlyContinue"
	}
	
	$record = Get-DnsServerResourceRecord @findArgs
	if ($Error.Count -gt 0) {
		if ($Error[0].CategoryInfo.Category -eq "ObjectNotFound")
		{
			$res = @{
				code = 404
				detail = "record not found"
			}
			Write-Output "$($res | ConvertTo-Json -Compress)"
			return
		}
		else {
				$res = @{
							code = 500
							detail = "$($Error[0].Exception.Message)"
				}
				Write-Output "$($res | ConvertTo-Json -Compress)"
				return
		}
	}
	
	$record = $record | Where-Object {
		$_.RecordData.IPv4Address.IPAddressToString -eq "{{.Address}}"
	}

	if ($null -eq $record) {
		$res = @{
			code = 404
			detail = "record not found"
		}
		Write-Output "$($res | ConvertTo-Json -Compress)"
		return
	}

	$newRecord = $record.Clone()
	if ({{.TTL}} -gt 0) {
		$newRecord.TimeToLive = [System.TimeSpan]::FromSeconds({{.TTL}})
	}
	if (![string]::IsNullOrEmpty("{{.NewAddress}}")) {
		$newRecord.RecordData.IPv4Address = [ipaddress]"{{.NewAddress}}"
	}

	$updateArgs = @{
		NewInputObject = $newRecord
		OldInputObject = $record
		ComputerName   = "{{.DnsServer}}"
		ZoneName       = "{{.ZoneName}}"
		Confirm        = $false
		ErrorAction    = "SilentlyContinue"
	}
	
	Set-DnsServerResourceRecord @updateArgs
	if ($Errors.Count -gt 0) {
		$res = @{
					code = 500
					detail = "$($Error[0].Exception.Message)"
			}
		Write-Output "$($res | ConvertTo-Json -Compress)"
		return
	}

	$records = @()
	$records += @{
		type = "A"
		name = $newRecord.HostName
		zone = "{{.ZoneName}}"
		data = $newRecord.RecordData.IPv4Address.IPAddressToString
		ttl  = $newRecord.TimeToLive.TotalSeconds
	}

	$res = @{
			code = 200
			detail = "record updated"
			records = $records
	}
	
	Write-Output "$($res | ConvertTo-Json -Compress -Depth 5)"
	`

	deleteARecordScript = `
	Import-Module DNSServer

	$findArgs = @{
		Name         = "{{.Name}}"
		ComputerName = "{{.DnsServer}}"
		ZoneName     = "{{.ZoneName}}"
		RRType       = "A"
		ErrorAction  = "SilentlyContinue"
	}
		
	$record = Get-DnsServerResourceRecord @findArgs
	if ($Error.Count -gt 0) {
		if ($Error[0].CategoryInfo.Category -eq "ObjectNotFound")
		{
			$res = @{
							code = 404
							detail = "record not found"
					}
			Write-Output "$($res | ConvertTo-Json -Compress)"
			return
		}
		else {
				$res = @{
							code = 500
							detail = "$($Error[0].Exception.Message)"
					}
				Write-Output "$($res | ConvertTo-Json -Compress)"
				return
		}
	}

	$record = $record | Where-Object {
		$_.RecordData.IPv4Address.IPAddressToString -eq "{{.Address}}"
	}

	if ($null -eq $record) {
		$res = @{
			code = 404
			detail = "record not found"
		}
		Write-Output "$($res | ConvertTo-Json -Compress)"
		return
	}

	$deleteArgs = @{
		ComputerName = "{{.DnsServer}}"
		ZoneName     = "{{.ZoneName}}"
		ErrorAction  = "SilentlyContinue"
		Confirm      = $false
		Force        = $true
	}

	$record | Remove-DnsServerResourceRecord @deleteArgs
	if ($Errors.Count -gt 0) {
		$res = @{
					code = 500
					detail = "$($Error[0].Exception.Message)"
			}
		Write-Output "$($res | ConvertTo-Json -Compress)"
		return
	}
	
	$records = @()
	$records += @{
		type = "A"
		name = $record.HostName
		zone = "{{.ZoneName}}"
		data = $record.RecordData.IPv4Address.IPAddressToString
		ttl  = $record.TimeToLive.TotalSeconds
	}

	$res = @{
			code    = 200
			detail  = "record deleted"
			records = $records
	}

	Write-Output "$($res | ConvertTo-Json -Compress)"
	`
)
