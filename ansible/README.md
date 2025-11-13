# Ansible Deployment for dpkg-build-pg

This directory contains Ansible playbooks for deploying the dpkg-build-pg application to Debian-based servers.

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

### 3. Configure Variables

Edit `group_vars/all.yml` to set your desired configuration:

```yaml
# Path to your .deb package
deb_package_path: "../dpkg-build-pg_1.0.0_amd64.deb"

# Server settings
server_port: 8080
server_host: "0.0.0.0"
server_message: "Hello from dpkg-build-pg server!"
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

### Option 3: Use Custom Config Template

If you want full control over the configuration file:

1. Uncomment the `custom_config_template` variable in `group_vars/all.yml`:
   ```yaml
   custom_config_template: "templates/config.yaml.j2"
   ```

2. The playbook will use the Jinja2 template instead of updating the existing config file

## Advanced Usage

### Override Variables at Runtime

```bash
# Override specific variables
ansible-playbook -i inventory playbook.yml \
  -e "server_port=9090" \
  -e "server_message='Custom message'"

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
4. **Install Package** - Installs dpkg-build-pg using apt
5. **Configure Application** - Updates configuration settings
6. **Enable Service** - Enables and starts the systemd service
7. **Verify Deployment** - Checks that the service is responding
8. **Cleanup** - Removes temporary files

## Reload Configuration

After the application is deployed, you can reload configuration without restarting:

```bash
ansible dpkg_servers -i inventory -b -m systemd -a "name=dpkg-build-pg state=reloaded"
```

Or use the systemctl command directly on the server:

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
├── group_vars/
│   └── all.yml                  # Default variables
└── templates/
    └── config.yaml.j2          # Configuration template (optional)
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