local fs = require "nixio.fs"
local sys = require "luci.sys"
local util = require "luci.util"

local m, s, o

local config_path = "/etc/config/koiddns"
local config_content = fs.readfile(config_path) or ""

m = Map("koiddns", translate("KoidDNS"),
    translate("Dynamic DNS Client for Aliyun and Tencent Cloud. Edit YAML below; use \"Validate\" to check format before saving."))

s = m:section(TypedSection, "yaml_config", translate("Configuration"))
s.anonymous = true
s.addremove = false

o = s:option(TextValue, "_yaml_content", translate("YAML"),
    translate("Configuration file content. Must be valid YAML with main, providers, and domains."))
o.rows = 22
o.wrap = "off"

function o.cfgvalue(self, section)
	return config_content
end

function o.write(self, section, value)
	local content = value:gsub("\r\n", "\n")
	local tmp = "/tmp/koiddns.luci.validate"
	fs.writefile(tmp, content)
	local ok = sys.call("/usr/bin/koiddns --config " .. tmp .. " --validate >/dev/null 2>&1") == 0
	fs.unlink(tmp)
	if not ok then
		return nil, translate("Format check failed. Fix YAML or use \"Validate format\" to see details.")
	end
	return fs.writefile(config_path, content)
end

s = m:section(SimpleSection)
s.template = "koiddns/apply"

return m
