# flwlnet-pat

Florida Winlink Net registry files for Pat P2P hosts. 

Create an RMS Express compatible global folder registry file for weekly submission to w4akh.net

## Interface

![main window](assets/screenshot.png)

## Usage

Build the code using the source, for your local system, or use the built flwlnetpat-linux-amd64 binary release.

- The default pat inbox will load automatically on start. 
  - If the default location is incorrect use the Load Mailbox button. 
- Select the messages to include
- Click Write Registry

## Make Targets

* `<none>` or `bin`: build for local os
* `tools`: install format, lint, and fyne build tools
* `lint`: run linter
* `format`: run format
* `release`: run fyne release build process
* `podman-init-<os>-<arch>`: create podman cross compiler image
* `podman-build-<os>-<arch>`: cross compile using podman image
