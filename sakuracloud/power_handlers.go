package sakuracloud

import (
	"context"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

const serverPowerAPILockKey = "sakuracloud_server.power.%d.lock"

func bootServerSync(ctx context.Context, client *APIClient, zone string, id types.ID) error {
	serverOp := sacloud.NewServerOp(client)
	if err := bootServer(ctx, client, zone, id); err != nil {
		return err
	}
	waiter := sacloud.WaiterForUp(func() (interface{}, error) { return serverOp.Read(ctx, zone, id) })
	if _, err := waiter.WaitForState(ctx); err != nil {
		return err
	}
	return nil
}

func shutdownServerSync(ctx context.Context, client *APIClient, zone string, id types.ID) error {
	serverOp := sacloud.NewServerOp(client)
	if err := shutdownServer(ctx, client, zone, id); err != nil {
		return err
	}
	waiter := sacloud.WaiterForDown(func() (interface{}, error) { return serverOp.Read(ctx, zone, id) })
	if _, err := waiter.WaitForState(ctx); err != nil {
		return err
	}
	return nil
}

func bootServer(ctx context.Context, client *APIClient, zone string, id types.ID) error {
	serverOp := sacloud.NewServerOp(client)

	lockServerPowerState(id)
	defer unlockServerPowerState(id)

	if err := serverOp.Boot(ctx, zone, id); err != nil {
		return err
	}
	return nil
}

func shutdownServer(ctx context.Context, client *APIClient, zone string, id types.ID) error {
	serverOp := sacloud.NewServerOp(client)

	lockServerPowerState(id)
	defer unlockServerPowerState(id)

	if err := serverOp.Shutdown(ctx, zone, id, &sacloud.ShutdownOption{
		Force: true, // TODO 後で
	}); err != nil {
		return err
	}
	return nil

}

func lockServerPowerState(id types.ID) {
	sakuraMutexKV.Lock(getServerPowerAPILockKey(id.Int64()))
	sakuraMutexKV.Lock(serverAPILockKey)
}

func unlockServerPowerState(id types.ID) {
	sakuraMutexKV.Unlock(serverAPILockKey)
	sakuraMutexKV.Unlock(getServerPowerAPILockKey(id.Int64()))
}

func getServerPowerAPILockKey(id int64) string {
	return fmt.Sprintf(serverPowerAPILockKey, id)
}

func shutdownVPCRouterSync(ctx context.Context, client *APIClient, zone string, id types.ID) error {
	vrOp := sacloud.NewVPCRouterOp(client)

	if err := vrOp.Shutdown(ctx, zone, id, &sacloud.ShutdownOption{}); err != nil {
		return err
	}
	waiter := sacloud.WaiterForDown(func() (interface{}, error) { return vrOp.Read(ctx, zone, id) })
	if _, err := waiter.WaitForState(ctx); err != nil {
		return err
	}
	return nil
}
