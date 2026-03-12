-- LuCI controller for KoidDNS (compatible with LuCI in OpenWRT 21.02+)
module("luci.controller.koiddns", package.seeall)

function index()
	if not nixio.fs.access("/usr/bin/koiddns") then
		return
	end

	entry({"admin", "services", "koiddns"}, firstchild(), _("KoidDNS"), 60).dependent = false

	entry({"admin", "services", "koiddns", "general"},
		cbi("koiddns/general"), _("General Settings"), 1)

	entry({"admin", "services", "koiddns", "providers"},
		cbi("koiddns/providers"), _("Providers"), 2)

	entry({"admin", "services", "koiddns", "domains"},
		cbi("koiddns/domains"), _("Domains"), 3)

	entry({"admin", "services", "koiddns", "status"},
		call("action_status"), _("Service"), 4)

	entry({"admin", "services", "koiddns", "validate"}, call("action_validate"))
end

function action_validate()
	local fs = require "nixio.fs"
	local util = require "luci.util"
	local tmpout = "/tmp/koiddns.validate.out"
	os.execute("/usr/bin/koiddns --config /etc/config/koiddns --validate >" .. tmpout .. " 2>&1")
	local out = fs.readfile(tmpout) or ""
	if fs.unlink then fs.unlink(tmpout) end
	local ok = (out:match("Config is valid") or out:match("valid")) and not out:match("failed")
	if ok and (out == "" or out:match("^%s*$")) then out = translate("Config is valid") end
	luci.template.render("koiddns/validate_result", { ok = ok, message = util.pcdata(out) })
end

function action_status()
	local http = require "luci.http"
	local sys = require "luci.sys"
	local verb = http.formvalue("verb")
	if verb == "start" or verb == "stop" or verb == "restart" then
		sys.init[verb]("koiddns")
		http.redirect(http.build_url("admin", "services", "koiddns", "status"))
		return
	end
	luci.template.render("koiddns/status", {
		running = sys.init.running("koiddns"),
		token = luci.dispatcher.context.authkey,
	})
end
