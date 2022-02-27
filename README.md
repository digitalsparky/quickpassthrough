# VFIO-enabler
A project to simplify setting up GPU passthrough for QuickEMU and libvirt


## TODO
* ~~Everything~~
* ~~Show general warning to user and inform about making a backup and general expectations~~
* Detect if user has an amd or intel CPU
* ~~Tell user to enable IOMMU (VT-d/AMD-v) on their motherboard and bootloader~~
* ~~Integrate ls-iommu and locate graphic cards (and detect if they are in their own IOMMU group)~~
* Enable and configure vfio modules
* Fetch the ID for the GPUs and generate the correct kernel arguments for grub and systemd-boot
* A menu system (i will need help by some bash wizards for this)
* Blacklist drivers? (some cards require blacklisting as softdep is not enough)
* Install vendor_reset kernel module? (maybe far future)
* Setup looking-glass? (far future maybe)