# Ansible Deployment for dpkg-build-pg

This directory contains Ansible playbooks for deploying the dpkg-build-pg application to Debian-based servers.

**Important**: This playbook only deploys the .deb package. Configuration is managed through the `config.production.yaml` file that is included in the .deb package itself. To change the configuration, you must update `config.production.yaml`, rebuild the .deb package, and redeploy.

## Prerequisites

- Ansible 2.9 or higher installed on your control machine
- SSH access to target servers
- Target servers running Debian/Ubuntu
- Sudo privileges on target servers
- Pre-built .deb package or URL to download it

## Quick Start

### 1. Install Ansible (if not already installed)

```bash
# Using Make (recommended)
make install-deps

# Or manually:
# On macOS
brew install ansible

# On Ubuntu/Debian
sudo apt update
sudo apt install ansible

# On other systems using pip
pip install ansible
```

### 2. Prepare Your Inventory

Copy the example inventory file and customize it for your environment:

```bash
cd ansible
cp inventory.example inventory
```

Edit `inventory` and add your server details:

```ini
[dpkg_servers]
your-server ansible_host=192.168.1.10 ansible_user=ubuntu ansible_ssh_private_key_file=~/.ssh/id_rsa
```

### 3. Configure Package Location

Edit `group_vars/all.yml` to specify your .deb package location:

```yaml
# Path to your .deb package
deb_package_path: "../dpkg-build-pg_1.0.0_amd64.deb"
```

### 4. Test Connection

Verify Ansible can connect to your servers:

```bash
# Using Make
make ping

# Or directly
ansible -i inventory dpkg_servers -m ping
```

### 5. Deploy the Application

Run the playbook:

```bash
# Using Make (recommended)
make deploy

# Or directly
ansible-playbook -i inventory playbook.yml
```

## Using Makefile Commands

The Makefile provides convenient shortcuts for common operations:

```bash
# View all available commands
make help

# Deploy the application
make deploy

# Dry-run (check what would change)
make deploy-check

# Check service status
make status

# Restart the service
make restart

# Reload configuration without restart
make reload

# Stop/start the service
make stop
make start

# View recent logs
make logs

# Test connection to servers
make ping

# Limit to specific hosts
make deploy LIMIT=prod-server-1
make status LIMIT=staging

# Enable verbose output
make deploy VERBOSE=-vvv
```

## Deployment Options

### Option 1: Deploy from Local .deb File

Set in `group_vars/all.yml`:

```yaml
deb_package_path: "../dpkg-build-pg_1.0.0_amd64.deb"
# deb_package_url: ""  # Comment this out or leave empty
```

### Option 2: Deploy from URL

Set in `group_vars/all.yml`:

```yaml
# deb_package_path: ""  # Comment this out or leave empty
deb_package_url: "https://github.com/your-org/dpkg-build-pg/releases/download/v1.0.0/dpkg-build-pg_1.0.0_amd64.deb"
```

## Configuration Management

Configuration is embedded in the .deb package via `config.production.yaml`. To update the configuration:

1. Edit `config.production.yaml` in your project root
2. Build a new .deb package
3. Deploy the new package using Ansible

The application will automatically use `/etc/dpkg-build-pg/config.production.yaml` if it exists.

## Advanced Usage

### Deploy Specific Package Version

```bash
# Deploy to specific hosts
ansible-playbook -i inventory playbook.yml --limit prod-server-1

# Use a different .deb package
ansible-playbook -i inventory playbook.yml \
  -e "deb_package_path=/path/to/dpkg-build-pg_2.0.0_amd64.deb"
```

### Deploy to Different Environments

Create environment-specific variable files:

```bash
# Create environment files
mkdir -p group_vars/staging group_vars/production

# Edit staging variables
vi group_vars/staging/vars.yml
```

Then run with the appropriate inventory:

```bash
ansible-playbook -i inventory playbook.yml --limit staging
```

### Dry Run (Check Mode)

Preview changes without applying them:

```bash
ansible-playbook -i inventory playbook.yml --check
```

### Verbose Output

Get more detailed output:

```bash
ansible-playbook -i inventory playbook.yml -v   # Basic verbose
ansible-playbook -i inventory playbook.yml -vv  # More verbose
ansible-playbook -i inventory playbook.yml -vvv # Debug level
```

## Playbook Tasks

The playbook performs the following tasks:

1. **Install Dependencies** - Ensures required packages are installed
2. **Create Temporary Directory** - For staging the .deb package
3. **Copy/Download Package** - Gets the .deb file to the server
4. **Install Package** - Installs dpkg-build-pg using apt (includes config.production.yaml)
5. **Reload Systemd** - Reloads systemd daemon
6. **Enable Service** - Enables the service to start on boot
7. **Restart Service** - Restarts the service with new configuration
8. **Verify Deployment** - Checks that the service is responding
9. **Cleanup** - Removes temporary files

## Reload Configuration (After Manual Config Changes)

If you manually edit the config file on the server and want to reload without redeploying:

```bash
# Using make
make reload

# Or directly
ansible dpkg_servers -i inventory -b -m systemd -a "name=dpkg-build-pg state=reloaded"
```

Or use systemctl directly on the server:

```bash
sudo systemctl reload dpkg-build-pg
```

## Uninstall

To remove the application:

```bash
ansible dpkg_servers -i inventory -b -m apt -a "name=dpkg-build-pg state=absent purge=yes"
```

## Troubleshooting

### Connection Issues

```bash
# Test SSH connection
ssh -i ~/.ssh/id_rsa ubuntu@your-server

# Check Ansible can reach the host
ansible -i inventory dpkg_servers -m ping

# Test with verbose output
ansible -i inventory dpkg_servers -m ping -vvv
```

### Permission Issues

Make sure your user has sudo privileges:

```bash
# Test sudo access
ansible -i inventory dpkg_servers -b -m shell -a "whoami"
```

### Service Not Starting

Check service status on the target server:

```bash
sudo systemctl status dpkg-build-pg
sudo journalctl -u dpkg-build-pg -n 50
```

### View Application Logs

```bash
ansible dpkg_servers -i inventory -b -m shell -a "journalctl -u dpkg-build-pg -n 20"
```

## Directory Structure

```
ansible/
├── README.md                    # This file
├── Makefile                     # Make commands for easy deployment
├── playbook.yml                 # Main deployment playbook
├── manage.yml                   # Service management playbook
├── inventory.example            # Example inventory file
├── inventory                    # Your actual inventory (not in git)
└── group_vars/
    └── all.yml                  # Default variables
```

## Best Practices

1. **Version Control**: Keep `inventory` file out of version control (add to .gitignore) as it may contain sensitive information
2. **Use Vault**: For sensitive data, use Ansible Vault:
   ```bash
   ansible-vault encrypt group_vars/all.yml
   ansible-playbook -i inventory playbook.yml --ask-vault-pass
   ```
3. **Test First**: Always test in a staging environment before deploying to production
4. **Idempotency**: The playbook is idempotent - you can run it multiple times safely
5. **Tags**: Use tags for selective execution (can be added to tasks as needed)

## Integration with CI/CD

### GitHub Actions Example

```yaml
- name: Deploy with Ansible
  run: |
    cd ansible
    ansible-playbook -i inventory playbook.yml \
      -e "deb_package_path=../dpkg-build-pg_${{ github.ref_name }}_amd64.deb"
  env:
    ANSIBLE_HOST_KEY_CHECKING: False
```

## Support

For issues or questions:
- Check the main project README
- Review Ansible documentation: https://docs.ansible.com/
- Check server logs: `sudo journalctl -u dpkg-build-pg`