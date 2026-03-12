local fs = require "nixio.fs"
local util = require "luci.util"

local function escape_html(s)
	if not s or s == "" then return "" end
	return (s:gsub("&", "&amp;"):gsub("<", "&lt;"):gsub(">", "&gt;"):gsub('"', "&quot;"))
end

local config_path = "/etc/config/koiddns"
local config_content = fs.readfile(config_path) or ""

local m = Map("koiddns", translate("KoidDNS Domains"), translate("Current domains (read-only). Edit in General Settings."))
local s = m:section(TypedSection, "yaml_info", translate("YAML"))
s.anonymous = true
s.addremove = false
local o = s:option(DummyValue, "_yaml_content", translate("Content"))
o.rawhtml = true
o.cfgvalue = function(self, section)
	return escape_html(config_content):gsub("\n", "<br>")
end
return m
