// Auto-generated file by goversioninfo. Do not edit.
package main

import (
	"encoding/json"

	"github.com/josephspurrier/goversioninfo"
)

func unmarshalGoVersionInfo(b []byte) goversioninfo.VersionInfo {
	vi := goversioninfo.VersionInfo{}
	json.Unmarshal(b, &vi)
	return vi
}

var versionInfo = unmarshalGoVersionInfo([]byte(`{
	"FixedFileInfo":{
		"FileVersion": {
			"Major": 1,
			"Minor": 0,
			"Patch": 0,
			"Build": 0
		},
		"ProductVersion": {
			"Major": 1,
			"Minor": 0,
			"Patch": 0,
			"Build": 0
		},
		"FileFlagsMask": "3f",
		"FileFlags": "",
		"FileOS": "040004",
		"FileType": "01",
		"FileSubType": "00"
	},
	"StringFileInfo":{
		"Comments": "Backup Manager",
		"CompanyName": "Kovengo Software Limited",
		"FileDescription": "FPanel backup manager",
		"FileVersion": "v1.0.0.0",
		"InternalName": "backupmng.exe",
		"LegalCopyright": "Copyright (c) 2020 Kovengo Software Limited",
		"LegalTrademarks": "",
		"OriginalFilename": "backupmng.go",
		"PrivateBuild": "",
		"ProductName": "FPanel",
		"ProductVersion": "v1.0.0.0",
		"SpecialBuild": ""
	},
	"VarFileInfo":{
		"Translation": {
			"LangID": 1033,
			"CharsetID": 1200
		}
	}
}`))
