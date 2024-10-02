package main

import (
	"bufio"
	"strings"
	"testing"

	"github.com/newrelic/infra-integrations-sdk/v3/data/inventory"
)

var testGetModulesCorrectFormat = `Loaded Modules:
 http_module (static)
 access_compat_module (shared)
 actions_module (shared)
 alias_module (shared)
 core_module (static)
 so_module (static)
`
var testGetModulesDifferentLinesFormat = `Loaded Modules:
 core_module (static)
 so_module (static)
 http_module (static)
 random text
 alias_module (shared)
 :
`
var testGetVersionCorrectFormat = `Server version: Apache/2.4.6 (CentOS)
Server built:   Nov 14 2016 18:04:44
Server's Module Magic Number: 20120211:24
Server loaded:  APR 1.4.8, APR-UTIL 1.5.2
Compiled using: APR 1.4.8, APR-UTIL 1.5.2
Architecture:   64-bit
Server MPM:     prefork
  threaded:     no
    forked:     yes (variable process count)
Server compiled with....
 -D APR_HAS_SENDFILE
`
var testGetVersionDifferentLinesFormat = `Random text
Server built:   Nov 14 2016 18:04:44
Server version: Apache/2.4.6 (CentOS)
Server MPM:     prefork
  threaded:     no
    forked:     yes (variable process count)
 :
`
var testWrongLinesFormat = `
Random text
Random text
:
`
var testEmptyInput = ``
var testCustomerModules = `AH00548: NameVirtualHost has no effect and will be removed in the next release /etc/httpd/conf.d/sites/1111.gohealth.com.conf:1
AH00112: Warning: DocumentRoot [/var/www/html/bcbsm.gohealthinsurance.com] does not exist
Loaded Modules:
 core_module (static)
 so_module (static)
 http_module (static)
 access_compat_module (shared)
 actions_module (shared)
 alias_module (shared)
 allowmethods_module (shared)
 auth_basic_module (shared)
 auth_digest_module (shared)
 authn_anon_module (shared)
 authn_core_module (shared)
 authn_dbd_module (shared)
 authn_dbm_module (shared)
 authn_file_module (shared)
 authn_socache_module (shared)
 authz_core_module (shared)
 authz_dbd_module (shared)
 authz_dbm_module (shared)
 authz_groupfile_module (shared)
 authz_host_module (shared)
 authz_owner_module (shared)
 authz_user_module (shared)
 autoindex_module (shared)
 cache_module (shared)
 cache_disk_module (shared)
 data_module (shared)
 dbd_module (shared)
 deflate_module (shared)
 dir_module (shared)
 dumpio_module (shared)
 echo_module (shared)
 env_module (shared)
 expires_module (shared)
 ext_filter_module (shared)
 filter_module (shared)
 headers_module (shared)
 include_module (shared)
 info_module (shared)
 log_config_module (shared)
 logio_module (shared)
 mime_magic_module (shared)
 mime_module (shared)
 negotiation_module (shared)
 remoteip_module (shared)
 reqtimeout_module (shared)
 rewrite_module (shared)
 setenvif_module (shared)
 slotmem_plain_module (shared)
 slotmem_shm_module (shared)
 socache_dbm_module (shared)
 socache_memcache_module (shared)
 socache_shmcb_module (shared)
 status_module (shared)
 substitute_module (shared)
 suexec_module (shared)
 unique_id_module (shared)
 unixd_module (shared)
 version_module (shared)
 vhost_alias_module (shared)
 dav_module (shared)
 dav_fs_module (shared)
 dav_lock_module (shared)
 lua_module (shared)
 mpm_worker_module (shared)
 proxy_module (shared)
 lbmethod_bybusyness_module (shared)
 lbmethod_byrequests_module (shared)
 lbmethod_bytraffic_module (shared)
 lbmethod_heartbeat_module (shared)
 proxy_ajp_module (shared)
 proxy_balancer_module (shared)
 proxy_connect_module (shared)
 proxy_express_module (shared)
 proxy_fcgi_module (shared)
 proxy_fdpass_module (shared)
 proxy_ftp_module (shared)
 proxy_http_module (shared)
 proxy_scgi_module (shared)
 proxy_wstunnel_module (shared)
 ssl_module (shared)
 systemd_module (shared)
 cgid_module (shared)
 foo bar_module
`

func TestParseGetModuleCorrectFormat(t *testing.T) {
	i := inventory.New()

	err := getModules(bufio.NewReader(strings.NewReader(testGetModulesCorrectFormat)), i)

	if len(i.Items()) != 6 {
		t.Error()
	}
	if i.Items()["modules/access_compat"]["value"] != "enabled" {
		t.Error()
	}
	if i.Items()["modules/alias"]["value"] != "enabled" {
		t.Error()
	}
	if i.Items()["modules/actions"]["value"] != "enabled" {
		t.Error()
	}
	if i.Items()["modules/http"]["value"] != "enabled" {
		t.Error()
	}
	if i.Items()["modules/so"]["value"] != "enabled" {
		t.Error()
	}
	if i.Items()["modules/core"]["value"] != "enabled" {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}

func TestParseGetModuleDifferentLinesFormat(t *testing.T) {
	i := inventory.New()
	err := getModules(bufio.NewReader(strings.NewReader(testGetModulesDifferentLinesFormat)), i)

	if len(i.Items()) != 4 {
		t.Error()
	}
	if i.Items()["modules/alias"]["value"] != "enabled" {
		t.Error()
	}
	if i.Items()["modules/http"]["value"] != "enabled" {
		t.Error()
	}
	if i.Items()["modules/so"]["value"] != "enabled" {
		t.Error()
	}
	if i.Items()["modules/core"]["value"] != "enabled" {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}

func TestParseGetModulesWrongLinesFormat(t *testing.T) {
	i := inventory.New()
	err := getModules(bufio.NewReader(strings.NewReader(testWrongLinesFormat)), i)

	if len(i.Items()) != 0 {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}

func TestParseGetModulesEmptyInput(t *testing.T) {
	i := inventory.New()
	err := getModules(bufio.NewReader(strings.NewReader(testEmptyInput)), i)

	if len(i.Items()) != 0 {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}

func TestParseGetVersionCorrectFormat(t *testing.T) {
	i := inventory.New()
	err := getVersion(bufio.NewReader(strings.NewReader(testGetVersionCorrectFormat)), i)

	if len(i.Items()) != 1 {
		t.Error()
	}
	if i.Items()["version"]["value"] != "Apache/2.4.6 (CentOS)" {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}

func TestParseGetVersionDifferentLinesFormat(t *testing.T) {
	i := inventory.New()
	err := getVersion(bufio.NewReader(strings.NewReader(testGetVersionDifferentLinesFormat)), i)

	if len(i.Items()) != 1 {
		t.Error()
	}
	if i.Items()["version"]["value"] != "Apache/2.4.6 (CentOS)" {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}

func TestParseGetVersionWrongLinesFormat(t *testing.T) {
	i := inventory.New()
	err := getVersion(bufio.NewReader(strings.NewReader(testWrongLinesFormat)), i)

	if len(i.Items()) != 0 {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}

func TestParseGetVersionEmptyInput(t *testing.T) {
	i := inventory.New()
	err := getVersion(bufio.NewReader(strings.NewReader(testEmptyInput)), i)

	if len(i.Items()) != 0 {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}

func TestParseGetModulesFromCustomerOutput(t *testing.T) {
	i := inventory.New()
	err := getModules(bufio.NewReader(strings.NewReader(testCustomerModules)), i)

	if len(i.Items()) == 0 {
		t.Error()
	}
	if err != nil {
		t.Error()
	}
}
