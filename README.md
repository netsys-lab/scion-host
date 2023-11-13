# Scion-Host: Run a cross-platform SCION stack in the OVGU Education AS
This repository contains binaries and code to setup a SCION endhost inside the OVGU Education AS. It should serve as dummy to create a more sophisticated tooling later.

So far Linux (amd64) and Windows (amd64) are supported. MacOS (arm64/amd64) will be next.

## Windows
To run the endhost stack on Windows, just download the `scion-host.exe` from this repository. 

You need to open a `Commandline Terminal` as Administrator to run the endhost stack, e.g. by typing `cmd` into Windows search, right click on `Command Line` and use `Open as Administrator`. Now move to the folder where you downloaded the `scion-host.exe` e.g. via `cd C:\Users\User\Downloads\`. If you are in the correct folder, type `scion-host.exe` in the terminal and press `Enter`. 

You should see some logging output and the program should stay open.

SCION is installed in your `Program Files` directory, e.g. under `C:\Program Files\scion\windowsx64`. In this folder you can also find the `scion.exe` binary. To use this, either open a `Commandline Terminal` in this directory, or add this directory to your Path, following [this documentation](https://www.autodesk.com/support/technical/article/caas/sfdcarticles/sfdcarticles/Adding-folder-path-to-Windows-PATH-environment-variable.html).

You need to be in the OvGU university network to connect to the SCION network properly. This can be done either by connecting to OvGU WIFI or LAN being on the campus, or using the [OvGU VPN](https://www.urz.ovgu.de/vpn-path-204,616.html). In the end, you need an OvGU IP Address in the range of `141.44.xx.xx`.

You can test SCION connectivity by using `scion.exe`. You may need to set your local address properly: Run `arp -a` in your terminal and find your OvGU IP address (should start with `141.44.`).

Now you can show paths to a given destination (e.g. Demokritos): `scion.exe showpaths -r 71-2546 --local <your_ovgu_ip>` or ping a SCION host `scion.exe ping 71-2546,127.0.0.1 --local <your_ovgu_ip>`.

**There will be a Windows Installer for SCION available soon, which makes it easier to run it on windows. Also the output of the SCION binary will be improved.**

## Linux
**Note: You need to stop any `scionlab` installations before running the `scion-host` binary, since they use similar ports.**

To run the endhost stack on Windows, just download the `scion-host` binary from this repository. Run it as root via `sudo ./scion-host` and keep this binary running. You should see some logging output and the program should stay open. SCION is installed under `/etc/scion-host/linuxx64/`, in this folder, you can also find the `scion` binary. You can use the binary from this folder or [add it to your Path](https://phoenixnap.com/kb/linux-add-to-path).

You need to be in the OvGU university network to connect to the SCION network properly. This can be done either by connecting to OvGU WIFI or LAN being on the campus, or using the [OvGU VPN](https://www.urz.ovgu.de/vpn-path-204,616.html). In the end, you need an OvGU IP Address in the range of `141.44.xx.xx`.

You can test SCION connectivity by using the `scion` binary. You may need to set your local address properly: Run `ip a` in your terminal and find your OvGU IP address (should start with `141.44.`).

Now you can show paths to a given destination (e.g. Demokritos): `scion showpaths -r 71-2546 --local <your_ovgu_ip>` or ping a SCION host `scion ping 71-2546,127.0.0.1 --local <your_ovgu_ip>`.

## Notes
This repo contains a lot of binaries, which is really bad for git, needs to be changed later, but for now it's easier...