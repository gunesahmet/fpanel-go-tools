package Websrvmng

import (
	"github.com/alexflint/go-arg"
)

var args struct {
	// Apache General
	Timeout          string `arg:"--timeout"`
	KeepAlive        string `arg:"--keep-alive"`
	KeepAliveTimeout string `arg:"--keep-alive-timeout"`
	MaxKeepAliveReq  string `arg:"--max-keep-alive-req"`
	HostnameLookups  string `arg:"--hostname-lookups"`
	LogLevel         string `arg:"--loglevel"`

	// MPM Prefork
	StartServer            string `arg:"--start-server"`
	MinSpareServer         string `arg:"--min-spare-server"`
	MaxSpareServer         string `arg:"--max-clients"`
	MaxRequestWorkers      string `arg:"--max-req-per-child"`
	MaxConnectionsPerChild string `arg:"--max-req-per-child"`

	// MPM Prefork Worker
	WorkerStartServer            string `arg:"--worker-start-server"`
	MinSpareThreads              string `arg:"--worker-min-spare"`
	MaxSpareThreads              string `arg:"--worker-max-spare"`
	ThreadLimit                  string `arg:"--worker-thread-limit"`
	ThreadsPerChild              string `arg:"--worker-thread-per-child"`
	WorkerMaxRequestWorkers      string `arg:"--worker-max-req-worker"`
	WorkerMaxConnectionsPerChild string `arg:"--worker-max-conn-per-child"`

	// NGINX
	WorkerProcesses           string `arg:"--worker_processes"`
	WorkerConnections         string `arg:"--worker_connections"`
	TypesHashMaxSize          string `arg:"--types_hash_max_size"`
	ServerNamesHashBucketSize string `arg:"--server_names_hash_bucket_size"`
	MultiAccept               string `arg:"--multi_accept"`
	Sendfile                  string `arg:"--sendfile"`
	TcpNopush                 string `arg:"--tcp_nopush"`
	TcpNodelay                string `arg:"--tcp_nodelay"`
	ServerTokens              string `arg:"--server_tokens"`
	ServerNameInRedirect      string `arg:"--server_name_in_redirect"`
	Gzip                      string `arg:"--gzip"`
	GzipVary                  string `arg:"--gzip_vary"`
	GzipBuffers               string `arg:"--gzip_buffers"`
	GzipCompLevel             string `arg:"--gzip_comp_level"`
	GzipHttpVersion           string `arg:"--gzip_http_version"`

	ClientMaxBodySize      string `arg:"--client_max_body_size"`
	ClientBodyBufferSize   string `arg:"--client_body_buffer_size"`
	ProxyConnectTimeout    string `arg:"--proxy_connect_timeout"`
	ProxySendTimeout       string `arg:"--proxy_send_timeout"`
	ProxyReadTimeout       string `arg:"--proxy_read_timeout"`
	ProxySendLowat         string `arg:"--proxy_send_lowat"`
	ProxyBufferSize        string `arg:"--proxy_buffer_size"`
	ProxyBuffers           string `arg:"--proxy_buffers"`
	ProxyBusyBuffersSize   string `arg:"--proxy_busy_buffers_size"`
	ProxyTempFileWriteSize string `arg:"--proxy_temp_file_write_size"`

	FileExt      string `arg:"--file-ext"`
	FileExtOctet string `arg:"--file-ext-octet"`
	ModuleName   string `arg:"--mod-name"`
	ModuleStatus string `arg:"--mod-status"`
	ServerType   string `arg:"--server-type"`

	Init           bool `arg:"--init"`
	InitVhosts     bool `arg:"--init-vhosts"`
	ListModules    bool `arg:"--list-modules"`
	LoadModules    bool `arg:"--load-modules"`
	ListMimetypes  bool `arg:"--list-mimetypes"`
	AddMimetype    bool `arg:"--add-mimetype"`
	RemoveMimetype bool `arg:"--remove-mimetype"`
	ServerSetting  bool `arg:"--server-setting"`
	SaveSetting    bool `arg:"--save-setting"`
}

var Timeout string
var KeepAlive string
var KeepAliveTimeout string
var HostnameLookups string
var LogLevel string
var MaxKeepAliveReq string

// MPM Prefork
var StartServer string
var MinSpareServer string
var MaxSpareServer string
var MaxRequestWorkers string
var MaxConnectionsPerChild string

// MPM Prefork Worker
var WorkerStartServer string
var MinSpareThreads string
var MaxSpareThreads string
var ThreadLimit string
var ThreadsPerChild string
var WorkerMaxRequestWorkers string
var WorkerMaxConnectionsPerChild string

// NGINX
var WorkerProcesses string
var WorkerConnections string
var TypesHashMaxSize string
var ServerNamesHashBucketSize string
var MultiAccept string
var Sendfile string
var TcpNopush string
var TcpNodelay string
var ServerTokens string
var ServerNameInRedirect string
var Gzip string
var GzipVary string
var GzipBuffers string
var GzipProxied string
var GzipCompLevel string
var GzipHttpVersion string
var GzipTypes string

var ClientMaxBodySize string
var ClientBodyBufferSize string
var ProxyConnectTimeout string
var ProxySendTimeout string
var ProxyReadTimeout string
var ProxySendLowat string
var ProxyBufferSize string
var ProxyBuffers string
var ProxyBusyBuffersSize string
var ProxyTempFileWriteSize string

var FileExt string
var FileExtOctet string
var ModuleName string
var ModuleStatus string
var ServerType string

var MainProcess string

func ArgInit() {
	arg.MustParse(&args)

	Timeout = args.Timeout
	KeepAlive = args.KeepAlive
	KeepAliveTimeout = args.KeepAliveTimeout
	MaxKeepAliveReq = args.MaxKeepAliveReq
	HostnameLookups = args.HostnameLookups
	LogLevel = args.LogLevel

	// MPM Prefork
	StartServer = args.StartServer
	MinSpareServer = args.MinSpareServer
	MaxSpareServer = args.MaxSpareServer
	MaxRequestWorkers = args.MaxRequestWorkers
	MaxConnectionsPerChild = args.MaxConnectionsPerChild

	// MPM Prefork Worker
	WorkerStartServer = args.WorkerStartServer
	MinSpareThreads = args.MinSpareThreads
	MaxSpareThreads = args.MaxSpareThreads
	ThreadLimit = args.ThreadLimit
	ThreadsPerChild = args.ThreadsPerChild
	WorkerMaxRequestWorkers = args.WorkerMaxRequestWorkers
	WorkerMaxConnectionsPerChild = args.WorkerMaxConnectionsPerChild

	// NGINX
	WorkerProcesses = args.WorkerProcesses
	WorkerConnections = args.WorkerConnections
	TypesHashMaxSize = args.TypesHashMaxSize
	ServerNamesHashBucketSize = args.ServerNamesHashBucketSize
	MultiAccept = args.MultiAccept
	Sendfile = args.Sendfile
	TcpNopush = args.TcpNopush
	TcpNodelay = args.TcpNodelay
	ServerTokens = args.ServerTokens
	ServerNameInRedirect = args.ServerNameInRedirect
	Gzip = args.Gzip
	GzipVary = args.GzipVary
	GzipBuffers = args.GzipBuffers

	ClientMaxBodySize = args.ClientMaxBodySize
	ClientBodyBufferSize = args.ClientBodyBufferSize
	ProxyConnectTimeout = args.ProxyConnectTimeout
	ProxySendTimeout = args.ProxySendTimeout
	ProxyReadTimeout = args.ProxyReadTimeout
	ProxySendLowat = args.ProxySendLowat
	ProxyBufferSize = args.ProxyBufferSize
	ProxyBuffers = args.ProxyBuffers
	ProxyBusyBuffersSize = args.ProxyBusyBuffersSize
	ProxyTempFileWriteSize = args.ProxyTempFileWriteSize

	GzipCompLevel = args.GzipCompLevel
	GzipHttpVersion = args.GzipHttpVersion
	FileExt = args.FileExt
	FileExtOctet = args.FileExtOctet
	ModuleName = args.ModuleName
	ModuleStatus = args.ModuleStatus
	ServerType = args.ServerType

	if args.Init {
		MainProcess = "init"
	} else if args.InitVhosts {
		MainProcess = "init-vhosts"
	} else if args.ListModules {
		MainProcess = "list-modules"
	} else if args.LoadModules {
		MainProcess = "load-modules"
	} else if args.ListMimetypes {
		MainProcess = "list-mimetypes"
	} else if args.AddMimetype {
		MainProcess = "add-mimetype"
	} else if args.RemoveMimetype {
		MainProcess = "remove-mimetype"
	} else if args.ServerSetting {
		MainProcess = "server-setting"
	} else if args.SaveSetting {
		MainProcess = "save-setting"
	}
}
