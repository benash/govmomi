package host

import (
	"context"
	"flag"
	"fmt"

	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

type info struct {
	*flags.ClientFlag
	*flags.HostSystemFlag
}

func init() {
	cli.Register("gpu.host.info", &info{})
}

func (cmd *info) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)

	cmd.HostSystemFlag, ctx = flags.NewHostSystemFlag(ctx)
	cmd.HostSystemFlag.Register(ctx, f)
}

func (cmd *info) Description() string {
	return `Display GPU information for a host.

Examples:
  govc gpu.host.info -host hostname`
}

func (cmd *info) Process(ctx context.Context) error {
	if err := cmd.ClientFlag.Process(ctx); err != nil {
		return err
	}
	if err := cmd.HostSystemFlag.Process(ctx); err != nil {
		return err
	}
	return nil
}

// Support NVIDIA devices, and exclude virtual GPUs, which have a SubDeviceId of 0x0000
func isPhysicalGPU(device types.HostPciDevice) bool {
	return device.VendorId == 0x10de && device.SubDeviceId != 0x0000
}

func printGPUInfo(gpu types.HostPciDevice) {
	fmt.Printf("PCI ID: %s\n", gpu.Id)
	fmt.Printf("  Device Name: %s\n", gpu.DeviceName)
	fmt.Printf("  Vendor Name: %s\n", gpu.VendorName)
	fmt.Printf("  Device ID: 0x%04x\n", gpu.DeviceId)
	fmt.Printf("  Vendor ID: 0x%04x\n", gpu.VendorId)
	fmt.Printf("  SubVendor ID: 0x%04x\n", gpu.SubVendorId)
	fmt.Printf("  SubDevice ID: 0x%04x\n", gpu.SubDeviceId)
	fmt.Printf("  Class ID: 0x%04x\n", gpu.ClassId)
	fmt.Printf("  Bus: 0x%02x, Slot: 0x%02x, Function: 0x%02x\n", gpu.Bus, gpu.Slot, gpu.Function)
	if gpu.ParentBridge != "" {
		fmt.Printf("  Parent Bridge: %s\n", gpu.ParentBridge)
	}
}

func (cmd *info) Run(ctx context.Context, f *flag.FlagSet) error {
	host, err := cmd.HostSystem()
	if err != nil {
		return err
	}

	if host == nil {
		return flag.ErrHelp
	}

	var h mo.HostSystem
	pc := property.DefaultCollector(host.Client())
	err = pc.RetrieveOne(ctx, host.Reference(), []string{"hardware"}, &h)
	if err != nil {
		return err
	}

	for _, device := range h.Hardware.PciDevice {
		if isPhysicalGPU(device) {
			printGPUInfo(device)
		}
	}

	return nil
}
