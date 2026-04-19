FROM openwrt/rootfs:x86_64-22.03.3@sha256:42bfe5aded99ddfbde5b9ba4c91d0826069d34eef6b537a6a510a5c4edf9147b

RUN mkdir /var/lock
RUN opkg update && opkg install \
    # Install curl so we can make a healthcheck
    # wget is installed, but it's hard to use for a health check.
    curl \
    # Install LuCI JSON-RPC packages.
    # See https://github.com/openwrt/luci/wiki/JsonRpcHowTo#basics
    luci-compat \
    luci-lib-ipkg \
    luci-mod-rpc \
    # Install LuCI (and HTTPS support)
    # This is entirely for debugging/diagnosis purposes.
    luci \
    luci-ssl
RUN /etc/init.d/uhttpd restart

HEALTHCHECK --interval=5s CMD curl \
    --data '{"id": 1, "method": "login", "params": ["root", ""]}' \
    --fail \
    --no-progress-meter \
    http://localhost/cgi-bin/luci/rpc/auth

EXPOSE 80 22

CMD ["/sbin/init"]
