# -*- coding: utf-8 -*-

INSTALL=install -v
UNINSTALL=rm -v -f

# BEL_BINDIR=../../bin
BEL_BINDIR=${GOPREFIX}/bin
BEL_CNFDIR=`pwd`

BINDIR=/usr/bin
CNFDIR=/etc/beluganos

SNMPD_CONF=/etc/snmp/snmpd.conf
SNMPD_PASS="/usr/bin/fibssnmp -v"
SNMPD_OID_LIST=`./snmpd_install.sh oid-list ./fibssnmp.yaml`
SNMPD_AGENT_PORT=8161

#
# services
#
.PHONY: install-service uninstall-service
install-service: install-beluganos install-fibc install-fibs install-govsw

uninstall-service: uninstall-beluganos uninstall-fibc uninstall-fibs uninstall-govsw


#
# beluganos
#
.PHONY: install-beluganos uninstall-beluganos
install-beluganos:
	@echo "install beluganos common service..."
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/beluganos.service /etc/systemd/system/
	@sudo systemctl daemon-reload
	@sudo systemctl enable beluganos.service
	@sudo systemctl start beluganos.service

uninstall-beluganos:
	@echo "uninstall beluganos common service..."
	@sudo systemctl disable beluganos.service || true
	@$(UNINSTALL) /etc/systemd/system/beluganos.service
	@sudo systemctl daemon-reload

#
# fibs
#
.PHONY: install-fibs uninstall-fibs
install-fibs: install-snmpproxyd install-fibssnmp install-fibsd register-snmp-oid set-snmpd-port

uninstall-fibs: uninstall-snmpproxyd uninstall-fibssnmp uninstall-fibsd unregister-snmp-oid unset-snmpd-port

#
# fibssnmp
#
.PHONY: install-fibssnmp uninstall-fibssnmp
install-fibssnmp:
	@echo "install fibssnmp..."
	@$(INSTALL) -pm 755 ${BEL_BINDIR}/fibssnmp      ${BINDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/fibssnmp.yaml ${CNFDIR}/

uninstall-fibssnmp:
	@echo "uninstall fibssnmp..."
	@$(UNINSTALL) ${CNFDIR}/fibssnmp.yaml
	@$(UNINSTALL) ${BINDIR}/fibssnmp

register-snmp-oid:
	@echo "register oid to snmpd.conf"
	@for oid in $(SNMPD_OID_LIST) ; do \
		./snmpd_install.sh install ${SNMPD_CONF} $$oid "${SNMPD_PASS}"; \
	done

unregister-snmp-oid:
	@echo "unregister oid to snmpd.conf"
	@for oid in $(SNMPD_OID_LIST) ; do \
		./snmpd_install.sh uninstall ${SNMPD_CONF} $$oid "${SNMPD_PASS}"; \
	done

#
# fibsd
#
.PHONY: install-fibsd uninstall-fibsd
install-fibsd:
	@echo "install fibsd..."
	@$(INSTALL) -pm 755 ${BEL_BINDIR}/fibsd         ${BINDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/fibsd.conf    ${CNFDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/fibsd.service /etc/systemd/system/
	@sudo systemctl daemon-reload

uninstall-fibsd:
	@echo "uninstall fibcd..."
	@$(UNINSTALL) /etc/systemd/system/fibsd.service
	@$(UNINSTALL) ${CNFDIR}/fibsd.conf
	@$(UNINSTALL) ${BINDIR}/fibsd
	@sudo systemctl daemon-reload


#
# fibs/snmpproxyd
#
.PHONY: install-snmpproxyd uninstall-snmpproxyd
install-snmpproxyd:
	@echo "install snmpproxyd..."
	@$(INSTALL) -pm 755 ${BEL_BINDIR}/snmpproxyd              ${BINDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/snmpproxyd.conf         ${CNFDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/snmpproxyd.yaml         ${CNFDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/snmpproxyd-mib.service  /etc/systemd/system/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/snmpproxyd-trap.service /etc/systemd/system/
	@sudo systemctl daemon-reload

uninstall-snmpproxyd:
	@echo "uninstall snmpproxyd..."
	@$(UNINSTALL) /etc/systemd/system/snmpproxyd-trap.service
	@$(UNINSTALL) /etc/systemd/system/snmpproxyd-mib.service
	@$(UNINSTALL) ${CNFDIR}/snmpproxyd.yaml
	@$(UNINSTALL) ${CNFDIR}/snmpproxyd.conf
	@$(UNINSTALL) ${BINDIR}/snmpproxyd
	@sudo systemctl daemon-reload


#
# snmpd.conf
#
.PHONY: set-snmpd-port unset-snmpd-port
set-snmpd-port:
	@echo "set snmpd listen port..."
	./snmpd_install.sh set-agent-port ${SNMPD_CONF} ${SNMPD_AGENT_PORT}
	@echo "enable snmpd rocommunity..."
	./snmpd_install.sh enable-rocommunity ${SNMPD_CONF}

unset-snmpd-port:
	@echo "unset snmpd listen port..."
	./snmpd_install.sh unset-agent-port ${SNMPD_CONF} ${SNMPD_AGENT_PORT}
	@echo "disable snmpd rocommunity..."
	./snmpd_install.sh disable-rocommunity ${SNMPD_CONF}


#
# fibc
#
.PHONY: install-fibc uninstall-fibc
install-fibc: install-fibcd install-fibcd-go
uninstall-fibc: uninstall-fibcd uninstall-fibcd-go

#
# fibcd (python)
#
.PHONY: install-fibcd uninstall-fibcd
install-fibcd:
	@echo "install fibcd..."
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/fibcd.conf    ${CNFDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/fibcd.service /etc/systemd/system/
	@sudo systemctl daemon-reload

uninstall-fibcd:
	@echo "uninstall fibcd..."
	@$(UNINSTALL) /etc/systemd/system/fibcd.service
	@$(UNINSTALL) ${CNFDIR}/fibcd.conf
	@sudo systemctl daemon-reload

#
# fibcd (golang)
#
.PHONY: install-fibcd-go uninstall-fibcd-go
install-fibcd-go:
	@echo "install fibcd(go)..."
	@$(INSTALL) -pm 755 ${BEL_BINDIR}/fibcd            ${BINDIR}/
	@$(INSTALL) -pm 755 ${BEL_BINDIR}/fibcctl          ${BINDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/fibcd-go.conf    ${CNFDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/fibcd-go.service /etc/systemd/system/
	@sudo systemctl daemon-reload

uninstall-fibcd-go:
	@echo "ubinstall fibcd(go)..."
	@$(UNINSTALL) /etc/systemd/system/fibcd-go.service
	@$(UNINSTALL) ${CNFDIR}/fibcd-go.conf
	@$(UNINSTALL) ${BINDIR}/fibcctl
	@$(UNINSTALL) ${BINDIR}/fibcd
	@sudo systemctl daemon-reload


#
# govsw
#
.PHONY: install-govsw uninstall-govsw
install-govsw: install-govswd
uninstall-govsw: uninstall-govswd

#
# govswd
#
.PHONY: install-govswd uninstall-govswd
install-govswd:
	@echo "install govswd"
	@$(INSTALL) -pm 755 ${BEL_BINDIR}/govswd         ${BINDIR}/
	@$(INSTALL) -pm 755 ${BEL_BINDIR}/govswc         ${BINDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/govswd.conf    ${CNFDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/govswd.yaml    ${CNFDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/govswd.service /etc/systemd/system/
	@sudo systemctl daemon-reload

uninstall-govswd:
	@echo "uninstall govswd"
	@$(UNINSTALL) /etc/systemd/system/govswd.service
	@$(UNINSTALL) ${CNFDIR}/govswd.yaml
	@$(UNINSTALL) ${CNFDIR}/govswd.conf
	@$(UNINSTALL) ${BINDIR}/govswd
	@$(UNINSTALL) ${BINDIR}/govswc
	@sudo systemctl daemon-reload
