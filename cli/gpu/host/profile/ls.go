package profile

import (
	"context"
	"flag"
	"fmt"

	"github.com/vmware/govmomi/cli"
	"github.com/vmware/govmomi/cli/flags"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
)

type ls struct {
	*flags.ClientFlag
	*flags.HostSystemFlag
}

func init() {
	cli.Register("gpu.host.profile.ls", &ls{})
}

func (cmd *ls) Register(ctx context.Context, f *flag.FlagSet) {
	cmd.ClientFlag, ctx = flags.NewClientFlag(ctx)
	cmd.ClientFlag.Register(ctx, f)

	cmd.HostSystemFlag, ctx = flags.NewHostSystemFlag(ctx)
	cmd.HostSystemFlag.Register(ctx, f)
}

func (cmd *ls) Description() string {
	return `List available vGPU profiles on host.

Examples:
  govc gpu.host.profile.ls -host hostname`
}

func (cmd *ls) Process(ctx context.Context) error {
	if err := cmd.ClientFlag.Process(ctx); err != nil {
		return err
	}
	if err := cmd.HostSystemFlag.Process(ctx); err != nil {
		return err
	}
	return nil
}

func (cmd *ls) Run(ctx context.Context, f *flag.FlagSet) error {
	host, err := cmd.HostSystem()
	if err != nil {
		return err
	}

	if host == nil {
		return flag.ErrHelp
	}

	var o mo.HostSystem
	pc := property.DefaultCollector(host.Client())
	err = pc.RetrieveOne(ctx, host.Reference(), []string{"config.sharedPassthruGpuTypes"}, &o)
	if err != nil {
		return err
	}

	if o.Config == nil {
		return fmt.Errorf("failed to get host configuration")
	}

	if len(o.Config.SharedPassthruGpuTypes) == 0 {
		return fmt.Errorf("no vGPU profiles available on this host")
	}

	fmt.Println("Available vGPU profiles:")
	for _, profile := range o.Config.SharedPassthruGpuTypes {
		fmt.Printf("  %s\n", profile)
	}

	return nil
}
