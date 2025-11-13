# GitHub Actions Workflows

This directory contains GitHub Actions workflows for building and releasing the Debian package.

## Workflows

### 1. Build Debian Package (`build-deb.yml`)

**Triggers:**
- Push to `main` branch
- Pull requests to `main` branch
- Manual trigger (workflow_dispatch)

**What it does:**
- Sets up Go environment
- Installs Debian packaging tools
- Builds the `.deb` package
- Uploads the package as an artifact

**Artifacts:**
- `debian-package`: The built `.deb` file (retained for 30 days)
- `build-artifacts`: Build metadata files (retained for 7 days)

**Usage:**
After a workflow run, you can download the `.deb` file from the Actions tab under "Artifacts".

### 2. Release Debian Package (`release.yml`)

**Triggers:**
- Push of tags matching `v*` (e.g., `v1.0.0`, `v1.2.3`)
- Manual trigger (workflow_dispatch)

**What it does:**
- Builds the Debian package
- Generates SHA256 checksums
- Creates a GitHub Release
- Attaches the `.deb` file and checksums to the release

**Creating a Release:**

1. Update `debian/changelog` with the new version
2. Commit your changes
3. Create and push a tag:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```
4. The workflow will automatically create a release with the package

## Requirements

Both workflows require:
- Ubuntu latest runner
- Go 1.21 or later
- Debian packaging tools (automatically installed)

## Permissions

The release workflow requires `contents: write` permission to create releases. This is automatically granted by the workflow configuration.

## Manual Triggers

Both workflows can be manually triggered from the Actions tab in GitHub:
1. Go to "Actions" tab
2. Select the workflow
3. Click "Run workflow"
4. Choose the branch
5. Click "Run workflow"

## Troubleshooting

### Build fails on dependency installation
Ensure the `debian/control` file has correct build dependencies.

### Release creation fails
Check that the `GITHUB_TOKEN` has proper permissions in repository settings under Settings → Actions → General → Workflow permissions.

### Package version mismatch
Ensure `debian/changelog` has the correct version number that matches your tag.
