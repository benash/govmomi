package vm

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
	*flags.VirtualMachineFlag
}

func init() {
	cli.Register("gpu.vm.info", &info{})
}

func (cmd *info) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)

	cmd.VirtualMachineFlag, ctx = flags.NewVirtualMachineFlag(ctx)
	cmd.VirtualMachineFlag.Register(ctx, f)
}

func (cmd *info) Process(ctx context.Context) error {
	if err := cmd.ClientFlag.Process(ctx); err != nil {
		return err
	}
	if err := cmd.VirtualMachineFlag.Process(ctx); err != nil {
		return err
	}
	return nil
}

func (cmd *info) Description() string {
	return `Display GPU information for a VM.

Examples:
  govc gpu.vm.info -vm $vm`
}

func (cmd *info) Run(ctx context.Context, f *flag.FlagSet) error {
	c, err := cmd.Client()
	if err != nil {
		return err
	}

	vm, err := cmd.VirtualMachine()
	if err != nil {
		return err
	}

	var o mo.VirtualMachine
	pc := property.DefaultCollector(c)
	err = pc.RetrieveOne(ctx, vm.Reference(), []string{"name", "config.hardware", "runtime.powerState"}, &o)
	if err != nil {
		return err
	}

	if o.Config == nil {
		return fmt.Errorf("VM configuration not available")
	}

	gpuCount := 0
	for _, device := range o.Config.Hardware.Device {
		if pciDevice, ok := device.(*types.VirtualPCIPassthrough); ok {
			fmt.Printf("GPU %d:\n", gpuCount)
			if desc := pciDevice.DeviceInfo.GetDescription(); desc != nil {
				fmt.Printf("  Label: %s\n", desc.Label)
				fmt.Printf("  Summary: %s\n", desc.Summary)
			}
			fmt.Printf("  Numa Node: %d\n", pciDevice.NumaNode)
			fmt.Printf("  Key: %d\n", pciDevice.Key)

			gpuCount++
		}
	}

	return nil
}
