package Websrvmng

import (
	"Helpers"
	"encoding/base64"
	"fmt"
	"strings"
)

func showApacheSettings() map[string]string {
	_, serverTimeout := Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("cat %s | grep \"^Timeout \" | awk '{print $2}'", apacheServer.ApacheConfPath))
	serverTimeout = Helpers.Trim(serverTimeout)

	_, keepAlive := Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("cat %s | grep \"^KeepAlive \" | awk '{print $2}'", apacheServer.ApacheConfPath))
	keepAlive = Helpers.Trim(keepAlive)

	_, maxKeepAliveRequest := Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("cat %s | grep \"^MaxKeepAliveRequests \" | awk '{print $2}'", apacheServer.ApacheConfPath))
	maxKeepAliveRequest = Helpers.Trim(maxKeepAliveRequest)

	_, keepAliveTimeout := Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("cat %s | grep \"^KeepAliveTimeout \" | awk '{print $2}'", apacheServer.ApacheConfPath))
	keepAliveTimeout = Helpers.Trim(keepAliveTimeout)

	_, hostnameLookups := Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("cat %s | grep \"^HostnameLookups \" | awk '{print $2}'", apacheServer.ApacheConfPath))
	hostnameLookups = Helpers.Trim(hostnameLookups)

	_, logLevel := Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("cat %s | grep \"^LogLevel \" | awk '{print $2}'", apacheServer.ApacheConfPath))
	logLevel = Helpers.Trim(hostnameLookups)

	jsonData := make(map[string]string)
	jsonData["hostname_lookups"] = hostnameLookups
	jsonData["keep_alive"] = keepAlive
	jsonData["keep_alive_timeout"] = keepAliveTimeout
	jsonData["max_keep_alive_requests"] = maxKeepAliveRequest
	jsonData["server_timeout"] = serverTimeout
	jsonData["log_level"] = logLevel

	jsonData = getMpmPreforkSettings(jsonData)
	jsonData = getMpmPreforkWorkerSettings(jsonData)

	return jsonData
}

func getMpmPreforkSettings(jsonData map[string]string) map[string]string {
	preFork := Helpers.ReadFile(fmt.Sprintf("%s/mpm_prefork.conf", apacheServer.ModulesAvailabeConfDir))

	var preForkStartServer, preForkMinSpareServer, preForkMaxSpareServer, preForkMaxRequestWorkers, preForkMaxConnectionsPerChild string
	for _, line := range preFork {
		line = Helpers.Trim(line)
		if Helpers.StringIsExists("StartServers", line) {
			preForkStartServer = Helpers.Trim(Helpers.RemoveAll(line, "StartServers"))
		}
		if Helpers.StringIsExists("MinSpareServers", line) {
			preForkMinSpareServer = Helpers.Trim(Helpers.RemoveAll(line, "MinSpareServers"))
		}
		if Helpers.StringIsExists("MaxSpareServers", line) {
			preForkMaxSpareServer = Helpers.Trim(Helpers.RemoveAll(line, "MaxSpareServers"))
		}
		if Helpers.StringIsExists("MaxRequestWorkers", line) {
			preForkMaxRequestWorkers = Helpers.Trim(Helpers.RemoveAll(line, "MaxRequestWorkers"))
		}
		if Helpers.StringIsExists("MaxConnectionsPerChild", line) {
			preForkMaxConnectionsPerChild = Helpers.Trim(Helpers.RemoveAll(line, "MaxConnectionsPerChild"))
		}
	}

	jsonData["start_server"] = preForkStartServer
	jsonData["min_spare_servers"] = preForkMinSpareServer
	jsonData["max_spare_servers"] = preForkMaxSpareServer
	jsonData["max_request_workers"] = preForkMaxRequestWorkers
	jsonData["max_conn_per_child"] = preForkMaxConnectionsPerChild
	return jsonData
}

func getMpmPreforkWorkerSettings(jsonData map[string]string) map[string]string {
	preForkWorker := Helpers.ReadFile(fmt.Sprintf("%s/mpm_worker.conf", apacheServer.ModulesAvailabeConfDir))

	var preForkWorkerStartServer, preForkWorkerMinSpareThreads, preForkWorkerMaxSpareThreads, preForkWorkerThreadLimit,
	preForkWorkerThreadsPerChild, preForkWorkerMaxRequestWorkers, preForkWorkerMaxConnPerDhild string

	for _, line := range preForkWorker {
		line = Helpers.Trim(line)
		if Helpers.StringIsExists("StartServers", line) {
			preForkWorkerStartServer = Helpers.Trim(Helpers.RemoveAll(line, "StartServers"))
		}
		if Helpers.StringIsExists("MinSpareThreads", line) {
			preForkWorkerMinSpareThreads = Helpers.Trim(Helpers.RemoveAll(line, "MinSpareThreads"))
		}
		if Helpers.StringIsExists("MaxSpareThreads", line) {
			preForkWorkerMaxSpareThreads = Helpers.Trim(Helpers.RemoveAll(line, "MaxSpareThreads"))
		}
		if Helpers.StringIsExists("ThreadLimit", line) {
			preForkWorkerThreadLimit = Helpers.Trim(Helpers.RemoveAll(line, "ThreadLimit"))
		}
		if Helpers.StringIsExists("ThreadsPerChild", line) {
			preForkWorkerThreadsPerChild = Helpers.Trim(Helpers.RemoveAll(line, "ThreadsPerChild"))
		}
		if Helpers.StringIsExists("MaxRequestWorkers", line) {
			preForkWorkerMaxRequestWorkers = Helpers.Trim(Helpers.RemoveAll(line, "MaxRequestWorkers"))
		}
		if Helpers.StringIsExists("MaxConnectionsPerChild", line) {
			preForkWorkerMaxConnPerDhild = Helpers.Trim(Helpers.RemoveAll(line, "MaxConnectionsPerChild"))
		}
	}

	jsonData["worker_start_server"] = preForkWorkerStartServer
	jsonData["worker_min_spare_threads"] = preForkWorkerMinSpareThreads
	jsonData["worker_max_spare_threads"] = preForkWorkerMaxSpareThreads
	jsonData["worker_thread_limit"] = preForkWorkerThreadLimit
	jsonData["worker_threads_per_child"] = preForkWorkerThreadsPerChild
	jsonData["worker_max_request_workers"] = preForkWorkerMaxRequestWorkers
	jsonData["worker_max_conn_per_child"] = preForkWorkerMaxConnPerDhild
	return jsonData
}

func saveApacheSettings() {
	varConf := map[string]string{
		"Timeout":              Timeout,
		"KeepAlive":            KeepAlive,
		"MaxKeepAliveRequests": MaxKeepAliveReq,
		"KeepAliveTimeout":     KeepAliveTimeout,
		"HostnameLookups":      HostnameLookups,
		"LogLevel":             LogLevel,
	}
	// For MPM
	varsMpm := map[string]string{
		"StartServers":           StartServer,
		"MinSpareServers":        MinSpareServer,
		"MaxSpareServers":        MaxSpareServer,
		"MaxRequestWorkers":      MaxRequestWorkers,
		"MaxConnectionsPerChild": MaxConnectionsPerChild,
	}
	// For MPM Worker
	varsMpmWorkers := map[string]string{
		"StartServers":           WorkerStartServer,
		"MinSpareThreads":        MinSpareThreads,
		"MaxSpareThreads":        MaxSpareThreads,
		"ThreadLimit":            ThreadLimit,
		"ThreadsPerChild":        ThreadsPerChild,
		"MaxRequestWorkers":      WorkerMaxRequestWorkers,
		"MaxConnectionsPerChild": WorkerMaxConnectionsPerChild,
	}

	var apacheConf string
	var mpmConf string
	var mpmWorkerConf string
	if Helpers.Compare(Env.SystemType, "ubuntu") {
		apacheConf = Helpers.Sprintf(ApacheUbuntuConfTemplate, varConf)
		mpmConf = Helpers.Sprintf(ApacheMpmConfTemplate, varsMpm)
		mpmWorkerConf = Helpers.Sprintf(ApacheMpmWorkerConfTemplate, varsMpmWorkers)
	} else if Helpers.Compare(Env.SystemType, "debian") {
		apacheConf = Helpers.Sprintf(ApacheDebianConfTemplate, varConf)
		mpmConf = Helpers.Sprintf(ApacheMpmConfTemplate, varsMpm)
		mpmWorkerConf = Helpers.Sprintf(ApacheMpmWorkerConfTemplate, varsMpm)
	} else if Helpers.Compare(Env.SystemType, "centos") {
		apacheConf = Helpers.Sprintf(ApacheCentosConfTemplate, varConf)
		mpmConf = Helpers.Sprintf(ApacheMpmConfTemplate, varsMpm)
		mpmWorkerConf = Helpers.Sprintf(ApacheMpmWorkerConfTemplate, varsMpm)
	}
	_ = Helpers.WriteFile(apacheConf, Env.ApacheConf)
	_ = Helpers.WriteFile(mpmConf, fmt.Sprintf("%s/mods-available/mpm_prefork.conf", Env.ApacheDir))
	_ = Helpers.WriteFile(mpmWorkerConf, fmt.Sprintf("%s/mods-available/mpm_worker.conf", Env.ApacheDir))
}

func apacheRemoveMimeType() {
	data, _ := base64.StdEncoding.DecodeString(FileExt)
	ext := Helpers.ReplaceAll(string(data), "#", " ")
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i 's/^.*%s$//g' %s", ext, Env.MimeTypesConfig))
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i '/^$/d' %s", Env.MimeTypesConfig))
}

func apacheAddMimeType() {
	_ = Helpers.AppendFile(fmt.Sprintf("%s\t\t\t\t%s", FileExtOctet, FileExt), "/etc/mime.types", 644)
	_, _ = Helpers.Exec("/bin/bash", "-c", fmt.Sprintf("sed -i '/^$/d' %s", Env.MimeTypesConfig))
}

func apacheMimeTypes() map[string]string {
	mimeTypes := make(map[string]string)
	if Helpers.FileExists("/etc/mime.types") {
		mimeTypesContent := Helpers.ReadFile("/etc/mime.types")
		for _, v := range mimeTypesContent {
			if strings.Contains(v, "#") {
				continue
			}
			v = Helpers.RemoveSpaces(v)
			vs := Helpers.Split(v, " ")
			if len(vs) > 1 {
				mimeTypes[vs[0]] = Helpers.ReplaceAllOld(v, vs[0], "")
			} else {
				mimeTypes[vs[0]] = ""
			}
		}
	}
	return mimeTypes
}

func apacheLoadModule(module string, status string) {
	if strings.Compare(status, "off") == 0 {
		if Helpers.FileExists(fmt.Sprintf("%s/%s.conf", apacheServer.ModulesEnabledConfDir, module)) {
			_ = Helpers.RemoveFile(fmt.Sprintf("%s/%s.conf", apacheServer.ModulesEnabledConfDir, module))
		}
		if Helpers.FileExists(fmt.Sprintf("%s/%s.load", apacheServer.ModulesEnabledConfDir, module)) {
			_ = Helpers.RemoveFile(fmt.Sprintf("%s/%s.load", apacheServer.ModulesEnabledConfDir, module))
		}
	} else {
		if Helpers.FileExists(fmt.Sprintf("%s/%s.conf", apacheServer.ModulesAvailabeConfDir, module)) {
			Helpers.CreateSymLink(fmt.Sprintf("%s/%s.conf", apacheServer.ModulesAvailabeConfDir, module),
				fmt.Sprintf("%s/%s.conf", apacheServer.ModulesEnabledConfDir, module))
		}
		if Helpers.FileExists(fmt.Sprintf("%s/%s.load", apacheServer.ModulesAvailabeConfDir, module)) {
			Helpers.CreateSymLink(fmt.Sprintf("%s/%s.load", apacheServer.ModulesAvailabeConfDir, module),
				fmt.Sprintf("%s/%s.load", apacheServer.ModulesEnabledConfDir, module))
		}
	}
}

func apacheInit(rewriteEnabled bool) {
	if rewriteEnabled {
		Helpers.CreateSymLink(fmt.Sprintf("%s/rewrite.conf", apacheServer.ModulesAvailabeConfDir),
			fmt.Sprintf("%s/rewrite.conf", apacheServer.ModulesEnabledConfDir))
	} else {
		_ = Helpers.RemoveFile(fmt.Sprintf("%s/rewrite.conf", apacheServer.ModulesEnabledConfDir))
	}
}