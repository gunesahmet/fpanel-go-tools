package Websrvmng

import (
	"Helpers"
	"encoding/base64"
	"fmt"
	"strings"
)

func showNginxSettings() map[string]string {
	jsonData := make(map[string]string)
	nginxConf := Helpers.ReadFile(Env.NginxDefaultConf)

	jsonData["worker_processes"] = ""
	jsonData["worker_connections"] = ""
	jsonData["multi_accept"] = ""
	jsonData["sendfile"] = ""
	jsonData["tcp_nopush"] = ""
	jsonData["tcp_nodelay"] = ""
	jsonData["keepalive_timeout"] = ""
	jsonData["types_hash_max_size"] = ""
	jsonData["server_tokens"] = ""
	jsonData["server_names_hash_bucket_size"] = ""
	jsonData["server_name_in_redirect"] = ""
	jsonData["gzip"] = ""
	jsonData["gzip_disable"] = ""
	jsonData["gzip_vary"] = ""
	jsonData["gzip_comp_level"] = ""
	jsonData["gzip_buffers"] = ""
	jsonData["gzip_http_version"] = ""
	jsonData["gzip_types"] = ""

	for _, line := range nginxConf {
		if len(line) == 0 {
			continue
		}

		emptyLine := Helpers.ExtractString(line, "^#")
		if len(emptyLine) > 0 {
			continue
		}
		disabledLine := Helpers.RemoveAfter(line, "#")
		if len(disabledLine) == 0 {
			continue
		}

		if strings.Contains(line, "worker_processes") {
			jsonData["worker_processes"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "worker_processes "), ";"), "#")
		}
		if strings.Contains(line, "worker_connections") {
			jsonData["worker_connections"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "worker_connections "), ";"), "#")
		}
		if strings.Contains(line, "multi_accept") {
			jsonData["multi_accept"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "multi_accept "), ";"), "#")
		}
		if strings.Contains(line, "sendfile") {
			jsonData["sendfile"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "sendfile "), ";"), "#")
		}
		if strings.Contains(line, "tcp_nopush") {
			jsonData["tcp_nopush"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "tcp_nopush "), ";"), "#")
		}
		if strings.Contains(line, "tcp_nodelay") {
			jsonData["tcp_nodelay"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "tcp_nodelay "), ";"), "#")
		}
		if strings.Contains(line, "keepalive_timeout") {
			jsonData["keepalive_timeout"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "keepalive_timeout "), ";"), "#")
		}
		if strings.Contains(line, "types_hash_max_size") {
			jsonData["types_hash_max_size"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "types_hash_max_size "), ";"), "#")
		}
		if strings.Contains(line, "server_tokens") {
			jsonData["server_tokens"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "server_tokens "), ";"), "#")
		}
		if strings.Contains(line, "server_names_hash_bucket_size") {
			jsonData["server_names_hash_bucket_size"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "server_names_hash_bucket_size "), ";"), "#")
		}
		if strings.Contains(line, "server_name_in_redirect") {
			jsonData["server_name_in_redirect"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "server_name_in_redirect "), ";"), "#")
		}
		if strings.Contains(line, "gzip ") {
			jsonData["gzip"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "gzip "), ";"), "#")
		}
		if strings.Contains(line, "gzip_disable") {
			jsonData["gzip_disable"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "gzip_disable "), ";"), "#")
		}
		if strings.Contains(line, "gzip_vary") {
			jsonData["gzip_vary"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "gzip_vary "), ";"), "#")
		}
		if strings.Contains(line, "gzip_comp_level") {
			jsonData["gzip_comp_level"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "gzip_comp_level "), ";"), "#")
		}
		if strings.Contains(line, "gzip_buffers") {
			jsonData["gzip_buffers"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "gzip_buffers "), ";"), "#")
		}
		if strings.Contains(line, "gzip_http_version") {
			jsonData["gzip_http_version"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "gzip_http_version "), ";"), "#")
		}
		if strings.Contains(line, "gzip_types") {
			jsonData["gzip_types"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "gzip_types "), ";"), "#")
		}
		if strings.Contains(line, "worker_rlimit_nofile") {
			jsonData["worker_rlimit_nofile"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "worker_rlimit_nofile "), ";"), "#")
		}
		if strings.Contains(line, "index ") {
			jsonData["index"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "index "), ";"), "#")
		}
		if strings.Contains(line, "default_type ") {
			jsonData["default_type"] = Helpers.RemoveAfter(Helpers.RemoveAll(Helpers.RemoveAll(Helpers.RemoveSpaces(line), "default_type "), ";"), "#")
		}
	}

	return jsonData
}

func saveNginxSettings() {
	vars := map[string]string{
		"worker_processes":              WorkerProcesses,
		"worker_connections":            WorkerConnections,
		"multi_accept":                  MultiAccept,
		"sendfile":                      Sendfile,
		"tcp_nopush":                    TcpNopush,
		"tcp_nodelay":                   TcpNodelay,
		"types_hash_max_size":           TypesHashMaxSize,
		"server_tokens":                 ServerTokens,
		"keepalive_timeout":             KeepAliveTimeout,
		"server_names_hash_bucket_size": ServerNamesHashBucketSize,
		"server_name_in_redirect":       ServerNameInRedirect,
		"gzip":                          Gzip,
		"gzip_vary":                     GzipVary,
		"gzip_proxied":                  GzipProxied,
		"gzip_comp_level":               GzipCompLevel,
		"gzip_buffers":                  GzipBuffers,
		"gzip_http_version":             GzipHttpVersion,
		"gzip_types":                    GzipTypes,

		"client_max_body_size":       ClientMaxBodySize,
		"client_body_buffer_size":    ClientBodyBufferSize,
		"proxy_connect_timeout":      ProxyConnectTimeout,
		"proxy_send_timeout":         ProxySendTimeout,
		"proxy_read_timeout":         ProxyReadTimeout,
		"proxy_send_lowat":           ProxySendLowat,
		"proxy_buffer_size":          ProxyBufferSize,
		"proxy_buffers":              ProxyBuffers,
		"proxy_busy_buffers_size":    ProxyBusyBuffersSize,
		"proxy_temp_file_write_size": ProxyTempFileWriteSize,
	}
	nginxConf := Helpers.Sprintf(NginxApacheConfTemplate, vars)
	_ = Helpers.WriteFile(nginxConf, Env.NginxDefaultConf)
}

func nginxRemoveMimeType() {
	data, _ := base64.StdEncoding.DecodeString(FileExt)
	ext := Helpers.ReplaceAll(string(data), "#", " ")
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i 's/^.*%s$//g' %s/mime.types", ext, Env.NginxDir))
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i '/^$/d' %s/mime.types", Env.NginxDir))
}

func nginxAddMimeType() {
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i 's/}//g' %s/mime.types", Env.NginxDir))
	_ = Helpers.AppendFile(fmt.Sprintf("    %s\t\t\t\t%s;", FileExtOctet, FileExt), fmt.Sprintf("%s/mime.types", Env.NginxDir), 644)
	_ = Helpers.AppendFile("}", fmt.Sprintf("%s/mime.types", Env.NginxDir), 644)
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i '/^$/d' %s/mime.types", Env.NginxDir))
}

func nginxMimeTypes() map[string]string {
	mimeTypes := make(map[string]string)
	if Helpers.FileExists(fmt.Sprintf("%s/mime.types", Env.NginxDir)) {
		mimeTypesContent := Helpers.ReadFile(fmt.Sprintf("%s/mime.types", Env.NginxDir))
		for _, v := range mimeTypesContent {
			if strings.Contains(v, "#") {
				continue
			}
			if len(Helpers.RemoveSpaces(v)) == 0 {
				continue
			}
			if strings.Contains(v, "types {") {
				continue
			}
			if strings.Contains(v, "}") {
				continue
			}
			v = Helpers.RemoveSpaces(v)
			vs := Helpers.Split(v, " ")
			mimeTypes[Helpers.RemoveAll(vs[0], ";")] = Helpers.RemoveAll(Helpers.ReplaceAllOld(v, vs[0], ""), ";")
			//mimeTypes[vs[0]] = vs[1]
		}
	}
	return mimeTypes
}

func nginxLoadModule(module string, status string) {
	PRIORITY := "50"
	if strings.Compare(module, "mod-http-ndk") == 0 {
		PRIORITY = "10"
	}
	if strings.Compare(status, "off") == 0 {
		if Helpers.FileExists(fmt.Sprintf("%s/%s-%s.conf", nginxServer.ModulesEnabledConfDir, PRIORITY, module)) {
			_ = Helpers.RemoveFile(fmt.Sprintf("%s/%s-%s.conf", nginxServer.ModulesEnabledConfDir, PRIORITY, module))
		}
	} else {
		if Helpers.FileExists(fmt.Sprintf("%s/%s-%s.conf", nginxServer.ModulesAvailabeConfDir, PRIORITY, module)) {
			Helpers.CreateSymLink(fmt.Sprintf("%s/%s.conf", nginxServer.ModulesAvailabeConfDir, module),
				fmt.Sprintf("%s/%s-%s.conf", nginxServer.ModulesEnabledConfDir, PRIORITY, module))
		}
	}
}