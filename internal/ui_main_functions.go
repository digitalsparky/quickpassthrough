package internal

import (
	"regexp"

	"github.com/HikariKnight/quickpassthrough/internal/configs"
	"github.com/HikariKnight/quickpassthrough/pkg/fileio"
)

// This function processes the enter event
func (m *model) processSelection() bool {
	switch m.focused {
	case GPUS:
		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// Gets the IOMMU group of the selected item
		iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
		iommu_group := iommu_group_regex.FindString(selectedItem.(item).desc)

		// Add the gpu group to our model (this is so we can grab the vbios details later)
		m.gpu_group = iommu_group

		// Get all the gpu devices and related devices (same device id or in the same group)
		items := iommuList2ListItem(getIOMMU("-grr", "-i", m.gpu_group, "-F", "vendor:,prod_name,optional_revision:,device_id"))

		// Add the devices to the list
		m.lists[GPU_GROUP].SetItems(items)

		// Change focus to next index
		m.focused++

	case GPU_GROUP:
		// Generate the VBIOS dumper script once the user has selected a GPU
		generateVBIOSDumper(*m)

		// Change focus to the next view
		m.focused++

	case USB:
		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// Gets the IOMMU group of the selected item
		iommu_group_regex := regexp.MustCompile(`(\d{1,3})`)
		iommu_group := iommu_group_regex.FindString(selectedItem.(item).desc)

		// Get the USB controllers in the selected iommu group
		items := iommuList2ListItem(getIOMMU("-ur", "-i", iommu_group, "-F", "vendor:,prod_name,optional_revision:,device_id"))

		// Add the items to the list
		m.lists[USB_GROUP].SetItems(items)

		// Change focus to next index
		m.focused++

	case USB_GROUP:
		m.focused++

	case VBIOS:
		// This is just an OK Dialog
		m.focused++

	case VIDEO:
		// This is a YESNO Dialog
		// Gets the selected item
		selectedItem := m.lists[m.focused].SelectedItem()

		// If user selected yes then
		if selectedItem.(item).title == "YES" {
			// Add disable VFIO video to the config
			configs.DisableVFIOVideo(1)
		} else {
			// Add disable VFIO video to the config
			configs.DisableVFIOVideo(0)
		}

		// Get our config struct
		config := configs.GetConfig()

		// If we have files for modprobe
		if fileio.FileExist(config.Path.MODPROBE) {
			// Get the device ids for the selected gpu using ls-iommu
			gpu_IDs := getIOMMU("-gr", "-i", m.gpu_group, "--id")

			// Configure modprobe
			configs.Set_Modprobe(gpu_IDs)
		}

		// If we have a folder for dracut
		if fileio.FileExist(config.Path.DRACUT) {
			// Configure dracut
			configs.Set_Dracut()
		}

		// Go to the next view
		//m.focused++

		// Because we have no QuickEmu support yet, just skip USB Controller configuration
		m.focused = DONE

	case INTRO:
		// This is an OK Dialog
		// Create the config folder and the files related to this system
		configs.InitConfigs()

		// Go to the next view
		m.focused++

	case DONE:
		// Return true so that the application will exit nicely
		return true
	}

	// Return false as we are not done
	return false
}
