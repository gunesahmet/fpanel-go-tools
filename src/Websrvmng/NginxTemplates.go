package Websrvmng

var NginxApacheConfTemplate = `
user  nginx;
%{worker_processes}s;
include /etc/nginx/modules-enabled/*.conf;

error_log  /var/log/nginx/error.log notice;
pid        /var/run/nginx.pid;

events {
    %{worker_connections}s;
	%{multi_accept}s;
}

http {
		%{sendfile}s;
        %{tcp_nopush}s;
        %{tcp_nodelay}s;
        %{types_hash_max_size}s;
        %{server_tokens}s;
	    %{keepalive_timeout}s;

        %{server_names_hash_bucket_size}s;
        %{server_name_in_redirect}s;

        include /etc/nginx/mime.types;
        default_type application/octet-stream;

        ssl_protocols TLSv1 TLSv1.1 TLSv1.2 TLSv1.3; # Dropping SSLv3, ref: POODLE
        ssl_prefer_server_ciphers on;

    	access_log  /var/log/nginx/access.log  main;
    	log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
        	              '$status $body_bytes_sent "$http_referer" '
            	          '"$http_user_agent" "$http_x_forwarded_for"';
        %{gzip}s;

        %{gzip_vary}s;
        %{gzip_proxied}s;
        %{gzip_comp_level}s;
        %{gzip_buffers}s;
        %{gzip_http_version}s;
        %{gzip_types}s;

		%{client_max_body_size}s;
		%{client_body_buffer_size}s;
			
		%{proxy_connect_timeout}s;
		%{proxy_send_timeout}s;
		%{proxy_read_timeout}s;
		%{proxy_send_lowat}s;
		
		%{proxy_buffer_size}s;
		%{proxy_buffers}s;
		%{proxy_busy_buffers_size}s;
		%{proxy_temp_file_write_size}s;

        include /etc/nginx/conf.d/*.conf;
}
`