// Copyright 2016-2019 terraform-provider-sakuracloud authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sakuracloud

import (
	"context"
	"fmt"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

func shutdownServerSync(ctx context.Context, client *APIClient, zone string, id types.ID, force bool) error {
	serverOp := sacloud.NewServerOp(client)
	if err := serverOp.Shutdown(ctx, zone, id, &sacloud.ShutdownOption{Force: force}); err != nil {
		return err
	}
	waiter := sacloud.WaiterForDown(func() (interface{}, error) { return serverOp.Read(ctx, zone, id) })
	if _, err := waiter.WaitForState(ctx); err != nil {
		return err
	}
	return nil
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

func bootDatabaseSync(ctx context.Context, dbOp sacloud.DatabaseAPI, zone string, id types.ID) error {
	if err := dbOp.Boot(ctx, zone, id); err != nil {
		return fmt.Errorf("booting Database[%s] is failed: %s", id, err)
	}
	waiter := sacloud.WaiterForUp(func() (interface{}, error) {
		return dbOp.Read(ctx, zone, id)
	})
	if _, err := waiter.WaitForState(ctx); err != nil {
		return fmt.Errorf("waiting for Database[%s] up is failed: %s", id, err)
	}
	return nil
}

func shutdownDatabaseSync(ctx context.Context, dbOp sacloud.DatabaseAPI, zone string, id types.ID, forceShutdown bool) error {
	if err := dbOp.Shutdown(ctx, zone, id, &sacloud.ShutdownOption{Force: forceShutdown}); err != nil {
		return fmt.Errorf("stopping Database[%s] is failed: %s", id, err)
	}
	waiter := sacloud.WaiterForDown(func() (interface{}, error) {
		return dbOp.Read(ctx, zone, id)
	})
	if _, err := waiter.WaitForState(ctx); err != nil {
		return fmt.Errorf("waiting for Database[%s] down is failed: %s", id, err)
	}
	return nil
}

func shutdownLoadBalancerSync(ctx context.Context, lbOp sacloud.LoadBalancerAPI, zone string, id types.ID, forceShutdown bool) error {
	if err := lbOp.Shutdown(ctx, zone, id, &sacloud.ShutdownOption{Force: forceShutdown}); err != nil {
		return fmt.Errorf("stopping SakuraCloud LoadBalancer[%s] is failed: %s", id, err)
	}
	waiter := sacloud.WaiterForDown(func() (interface{}, error) {
		return lbOp.Read(ctx, zone, id)
	})
	if _, err := waiter.WaitForState(ctx); err != nil {
		return fmt.Errorf("stopping SakuraCloud LoadBalancer[%s] is failed: %s", id, err)
	}
	return nil
}

func shutdownNFSSync(ctx context.Context, nfsOp sacloud.NFSAPI, zone string, id types.ID, forceShutdown bool) error {
	if err := nfsOp.Shutdown(ctx, zone, id, &sacloud.ShutdownOption{Force: forceShutdown}); err != nil {
		return fmt.Errorf("stopping SakuraCloud NFS[%s] is failed: %s", id, err)
	}
	waiter := sacloud.WaiterForDown(func() (interface{}, error) {
		return nfsOp.Read(ctx, zone, id)
	})
	if _, err := waiter.WaitForState(ctx); err != nil {
		return fmt.Errorf("waiting for NFS[%s] down is failed: %s", id, err)
	}
	return nil
}
