package Websrvmng

var ApacheUbuntuConfTemplate = `Mutex file:${APACHE_LOCK_DIR} default
PidFile ${APACHE_PID_FILE}
Timeout %{Timeout}s
KeepAlive %{KeepAlive}s
MaxKeepAliveRequests %{MaxKeepAliveRequests}s
KeepAliveTimeout %{KeepAliveTimeout}s
User ${APACHE_RUN_USER}
Group ${APACHE_RUN_GROUP}
HostnameLookups %{HostnameLookups}s
ErrorLog ${APACHE_LOG_DIR}/error.log
LogLevel %{LogLevel}s

IncludeOptional mods-enabled/*.load
IncludeOptional mods-enabled/*.conf

Include ports.conf

<Directory />
        Options FollowSymLinks
        AllowOverride None
        Require all denied
</Directory>

<Directory /usr/share>
        AllowOverride None
        Require all granted
</Directory>

<Directory /var/www/>
        Options Indexes FollowSymLinks
        AllowOverride None
        Require all granted
</Directory>

AccessFileName .htaccess

<FilesMatch "^\.ht">
        Require all denied
</FilesMatch>

LogFormat "%v:%p %h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" vhost_combined
LogFormat "%h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" combined
LogFormat "%h %l %u %t \"%r\" %>s %O" common
LogFormat "%{Referer}i -> %U" referer
LogFormat "%{User-agent}i" agent

IncludeOptional conf.d/*.conf
`
var ApacheDebianConfTemplate = `Mutex file:${APACHE_LOCK_DIR} default
DefaultRuntimeDir ${APACHE_RUN_DIR}
PidFile ${APACHE_PID_FILE}
Timeout %{Timeout}s
KeepAlive %{KeepAlive}s
MaxKeepAliveRequests %{MaxKeepAliveRequests}s
KeepAliveTimeout %{KeepAliveTimeout}s
User ${APACHE_RUN_USER}
Group ${APACHE_RUN_GROUP}
HostnameLookups %{HostnameLookups}s
ErrorLog ${APACHE_LOG_DIR}/error.log
LogLevel %{LogLevel}s
IncludeOptional mods-enabled/*.load
IncludeOptional mods-enabled/*.conf
Include ports.conf
<Directory />
        Options FollowSymLinks
        AllowOverride None
        Require all denied
</Directory>

<Directory /usr/share>
        AllowOverride None
        Require all granted
</Directory>

<Directory /var/www/>
        Options Indexes FollowSymLinks
        AllowOverride None
        Require all granted
</Directory>

AccessFileName .htaccess
<FilesMatch "^\.ht">
        Require all denied
</FilesMatch>
LogFormat "%v:%p %h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" vhost_combined
LogFormat "%h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" combined
LogFormat "%h %l %u %t \"%r\" %>s %O" common
LogFormat "%{Referer}i -> %U" referer
LogFormat "%{User-agent}i" agent
IncludeOptional conf.d/*.conf
`
var ApacheCentosConfTemplate = `ServerRoot "/etc/httpd"
Listen 80
Include mods-enabled/*.load
Include mods-enabled/*.conf

User apache
Group apache
ServerAdmin root@localhost

<Directory />
    AllowOverride none
    Require all denied
</Directory>
<Directory "/var/www">
    AllowOverride None
    Require all granted
</Directory>
<IfModule dir_module>
    DirectoryIndex index.html
</IfModule>
<Files ".ht*">
    Require all denied
</Files>
ErrorLog "logs/error_log"
LogLevel warn
<IfModule log_config_module>
        LogFormat "%v:%p %h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" vhost_combined
        LogFormat "%h %l %u %t \"%r\" %>s %O \"%{Referer}i\" \"%{User-Agent}i\"" combined
        LogFormat "%h %l %u %t \"%r\" %>s %O" common
    <IfModule logio_module>
#      LogFormat "%h %l %u %t \"%r\" %>s %b \"%{Referer}i\" \"%{User-Agent}i\" %I %O" combinedio
    </IfModule>
    CustomLog "logs/access_log" combined
</IfModule>

<IfModule mime_module>
    TypesConfig /etc/mime.types
    AddType application/x-compress .Z
    AddType application/x-gzip .gz .tgz
    AddType text/html .shtml
    AddOutputFilter INCLUDES .shtml
</IfModule>
AddDefaultCharset UTF-8
<IfModule mime_magic_module>
    MIMEMagicFile conf/magic
</IfModule>
EnableSendfile on
IncludeOptional conf.d/*.conf
`


var ApacheMpmConfTemplate = `<IfModule mpm_prefork_module>
        StartServers           %{StartServers}s
        MinSpareServers        %{MinSpareServers}s
        MaxSpareServers        %{MaxSpareServers}s
        MaxRequestWorkers      %{MaxRequestWorkers}s
        MaxConnectionsPerChild %{MaxConnectionsPerChild}s
</IfModule>`

var ApacheMpmWorkerConfTemplate = `<IfModule mpm_worker_module>
        StartServers           %{StartServers}s
        MinSpareThreads        %{MinSpareThreads}s
        MaxSpareThreads        %{MaxSpareThreads}s
        ThreadLimit            %{ThreadLimit}s
        ThreadsPerChild        %{ThreadsPerChild}s
        MaxRequestWorkers      %{MaxRequestWorkers}s
        MaxConnectionsPerChild %{MaxConnectionsPerChild}s
</IfModule>`