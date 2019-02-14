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

#
# services
#
.PHONY: install-service uninstall-service
install-service: install-fibc install-fibs

uninstall-service: uninstall-fibc uninstall-fibs


#
# fibs
#
.PHONY: install-fibs uninstall-fibs
install-fibs: install-snmpproxyd install-fibssnmp install-fibsd register-snmp-oid

uninstall-fibs: uninstall-snmpproxyd uninstall-fibssnmp uninstall-fibsd unregister-snmp-oid


#
# fibssnmp
#
.PHONY: install-fibssnmp uninstall-fibssnmp
install-fibssnmp:
	@echo "install fibssnmp..."
	@$(INSTALL) -pm 755 ${BEL_BINDIR}/fibssnmp      ${BINDIR}/
	@$(INSTALL) -pm 644 ${BEL_CNFDIR}/fibssnmp.yaml ${CNFDIR}/
	@echo "*** edit snmpd.conf and add 'pass_psesist <oid> ${BINDIR}/fibssnmp' ***"

uninstall-fibssnmp:
	@echo "uninstall fibssnmp..."
	@$(UNINSTALL) ${CNFDIR}/fibssnmp.yaml
	@$(UNINSTALL) ${BINDIR}/fibssnmp
	@echo "*** edit snmpd.conf and remove 'pass_psesist ...' ***"

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
# fibc
#
.PHONY: install-fibc uninstall-fibc
install-fibc: install-fibcd

uninstall-fibc: uninstall-fibcd


#
# fibcd
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
