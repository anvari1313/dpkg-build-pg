.PHONY: help deploy status restart reload stop start check ping clean install-deps

# Default inventory file
INVENTORY ?= inventory
PLAYBOOK ?= playbook.yml

help:
	@echo "Ansible deployment commands for dpkg-build-pg"
	@echo ""
	@echo "Usage:"
	@echo "  make deploy         Deploy application to servers"
	@echo "  make deploy-check   Dry-run deployment (check mode)"
	@echo "  make status         Check service status"
	@echo "  make restart        Restart the service"
	@echo "  make reload         Reload configuration without restart"
	@echo "  make stop           Stop the service"
	@echo "  make start          Start the service"
	@echo "  make ping           Test connection to servers"
	@echo "  make check          Check playbook syntax"
	@echo "  make clean          Clean up temporary files"
	@echo "  make install-deps   Install Ansible and dependencies"
	@echo ""
	@echo "Variables:"
	@echo "  INVENTORY=<file>   Specify inventory file (default: inventory)"
	@echo "  LIMIT=<hosts>      Limit to specific hosts"
	@echo "  VERBOSE=-vvv       Enable verbose output"
	@echo ""
	@echo "Examples:"
	@echo "  make deploy LIMIT=prod-server-1"
	@echo "  make status VERBOSE=-v"
	@echo "  make deploy INVENTORY=staging-inventory"

install-deps:
	@echo "Installing Ansible..."
	@which ansible > /dev/null 2>&1 || pip3 install ansible
	@echo "Ansible version:"
	@ansible --version

ping:
	ansible -i $(INVENTORY) dpkg_servers -m ping $(VERBOSE)

check:
	ansible-playbook -i $(INVENTORY) $(PLAYBOOK) --syntax-check
	@echo "Playbook syntax is valid"

deploy: check
	ansible-playbook -i $(INVENTORY) $(PLAYBOOK) $(VERBOSE) $(if $(LIMIT),--limit $(LIMIT),)

deploy-check: check
	ansible-playbook -i $(INVENTORY) $(PLAYBOOK) --check --diff $(VERBOSE) $(if $(LIMIT),--limit $(LIMIT),)

status:
	ansible-playbook -i $(INVENTORY) manage.yml -e "service_action=status" $(VERBOSE) $(if $(LIMIT),--limit $(LIMIT),)

restart:
	ansible-playbook -i $(INVENTORY) manage.yml -e "service_action=restart" $(VERBOSE) $(if $(LIMIT),--limit $(LIMIT),)

reload:
	ansible-playbook -i $(INVENTORY) manage.yml -e "service_action=reload" $(VERBOSE) $(if $(LIMIT),--limit $(LIMIT),)

stop:
	ansible-playbook -i $(INVENTORY) manage.yml -e "service_action=stop" $(VERBOSE) $(if $(LIMIT),--limit $(LIMIT),)

start:
	ansible-playbook -i $(INVENTORY) manage.yml -e "service_action=start" $(VERBOSE) $(if $(LIMIT),--limit $(LIMIT),)

logs:
	ansible -i $(INVENTORY) dpkg_servers -b -m shell -a "journalctl -u dpkg-build-pg -n 50 --no-pager" $(if $(LIMIT),--limit $(LIMIT),)

clean:
	find . -name "*.retry" -delete
	@echo "Cleaned up temporary files"