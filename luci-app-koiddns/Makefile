include $(TOPDIR)/rules.mk

PKG_NAME:=luci-app-koiddns
PKG_VERSION:=1.1
PKG_RELEASE:=1
# LuCI/OpenWRT: tested with 21.02, 22.03, 23.05; procd init for OpenWRT 15.05+

PKG_MAINTAINER:=Your Name <your.email@example.com>

include $(INCLUDE_DIR)/package.mk

define Package/luci-app-koiddns
  SECTION:=luci
  CATEGORY:=LuCI
  SUBMENU:=3. Applications
  TITLE:=KoidDNS Dynamic DNS Client
  PKGARCH:=all
  DEPENDS:=+koiddns
endef

define Package/luci-app-koiddns/description
  LuCI support for KoidDNS, a Dynamic DNS client for Aliyun and Tencent Cloud.
  This package provides a web interface for managing the YAML configuration file.
  Compatible with OpenWRT 21.02+ and koiddns (--config only; daemon/pidfile/logfile
  are not used when run via init/procd).
endef

define Package/luci-app-koiddns/conffiles
/etc/config/koiddns
endef

define Build/Prepare
endef

define Build/Configure
endef

define Build/Compile
endef

define Package/luci-app-koiddns/install
	$(INSTALL_DIR) $(1)/usr/lib/lua/luci/controller
	$(INSTALL_DATA) ./luasrc/controller/koiddns.lua $(1)/usr/lib/lua/luci/controller/koiddns.lua

	$(INSTALL_DIR) $(1)/usr/lib/lua/luci/model/cbi/koiddns
	$(INSTALL_DATA) ./luasrc/model/cbi/koiddns/general.lua $(1)/usr/lib/lua/luci/model/cbi/koiddns/general.lua
	$(INSTALL_DATA) ./luasrc/model/cbi/koiddns/providers.lua $(1)/usr/lib/lua/luci/model/cbi/koiddns/providers.lua
	$(INSTALL_DATA) ./luasrc/model/cbi/koiddns/domains.lua $(1)/usr/lib/lua/luci/model/cbi/koiddns/domains.lua

	$(INSTALL_DIR) $(1)/usr/lib/lua/luci/view/koiddns
	$(INSTALL_DATA) ./luasrc/view/koiddns/apply.htm $(1)/usr/lib/lua/luci/view/koiddns/apply.htm
	$(INSTALL_DATA) ./luasrc/view/koiddns/validate_result.htm $(1)/usr/lib/lua/luci/view/koiddns/validate_result.htm

	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) ./root/etc/config/koiddns $(1)/etc/config/koiddns

	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) ./root/etc/init.d/koiddns $(1)/etc/init.d/koiddns
endef

$(eval $(call BuildPackage,luci-app-koiddns))