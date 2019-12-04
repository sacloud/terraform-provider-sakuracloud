package sakuracloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
	"sync"
	"time"
)

type deletionWaiterFindFunc func(context.Context, *APIClient, string, types.ID) (bool, error)

func waitForDeletionBySwitchID(ctx context.Context, client *APIClient, zone string, switchID types.ID) error {
	finder := []deletionWaiterFindFunc{
		findServerBySwitchID,
		findLoadBalancerBySwitchID,
		findVPCRouterBySwitchID,
		findDatabaseBySwitchID,
		findNFSBySwitchID,
		findMobileGatewayBySwitchID,
	}

	var wg sync.WaitGroup
	wg.Add(len(finder))

	errCh := make(chan error)
	compCh := make(chan struct{})

	for _, f := range finder {
		go func(f deletionWaiterFindFunc) {
			if err := waitForDeletionByFunc(ctx, client, zone, switchID, f); err != nil {
				errCh <- err
			}
			wg.Done()
		}(f)
	}

	go func() {
		wg.Wait()
		compCh <- struct{}{}
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case <-compCh:
			return nil
		case <-time.After(client.deletionWaiterTimeout):
			return errors.New("waiting deletion is failed: timeout")
		}
	}
}

func waitForDeletionByFunc(ctx context.Context, client *APIClient, zone string, switchID types.ID, f deletionWaiterFindFunc) error {
	t := time.NewTicker(client.deletionWaiterPollingInterval)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			isExists, err := f(ctx, client, zone, switchID)
			if err != nil {
				return err
			}
			if !isExists {
				return nil
			}

		case <-time.After(client.deletionWaiterTimeout):
			return errors.New("waiting deletion is failed: timeout")
		}
	}
}

func findServerBySwitchID(ctx context.Context, client *APIClient, zone string, id types.ID) (bool, error) {
	swOp := sacloud.NewSwitchOp(client)

	searched, err := swOp.GetServers(ctx, zone, id)
	if err != nil {
		return false, fmt.Errorf("finding server is failed: %s", err)
	}
	return searched.Count != 0, nil
}

func findVPCRouterBySwitchID(ctx context.Context, client *APIClient, zone string, id types.ID) (bool, error) {
	vrOp := sacloud.NewVPCRouterOp(client)

	searched, err := vrOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding VPCRouter is failed: %s", err)
	}

	for _, vpcRouter := range searched.VPCRouters {
		for _, iface := range vpcRouter.Interfaces {
			if iface.SwitchScope != types.Scopes.Shared && iface.SwitchID == id {
				return true, nil
			}
		}
	}
	return false, nil
}

func findLoadBalancerBySwitchID(ctx context.Context, client *APIClient, zone string, id types.ID) (bool, error) {
	lbOp := sacloud.NewLoadBalancerOp(client)

	searched, err := lbOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding LoadBalancer is failed: %s", err)
	}

	for _, lb := range searched.LoadBalancers {
		for _, iface := range lb.Interfaces {
			if iface.SwitchScope != types.Scopes.Shared && iface.SwitchID == id {
				return true, nil
			}
		}
	}
	return false, nil
}

func findDatabaseBySwitchID(ctx context.Context, client *APIClient, zone string, id types.ID) (bool, error) {
	dbOp := sacloud.NewDatabaseOp(client)

	searched, err := dbOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding Database is failed: %s", err)
	}

	for _, db := range searched.Databases {
		for _, iface := range db.Interfaces {
			if iface.SwitchScope != types.Scopes.Shared && iface.SwitchID == id {
				return true, nil
			}
		}
	}
	return false, nil
}

func findNFSBySwitchID(ctx context.Context, client *APIClient, zone string, id types.ID) (bool, error) {
	nfsOp := sacloud.NewNFSOp(client)

	searched, err := nfsOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding NFS is failed: %s", err)
	}

	for _, nfs := range searched.NFS {
		for _, iface := range nfs.Interfaces {
			if iface.SwitchScope != types.Scopes.Shared && iface.SwitchID == id {
				return true, nil
			}
		}
	}
	return false, nil
}

func findMobileGatewayBySwitchID(ctx context.Context, client *APIClient, zone string, id types.ID) (bool, error) {
	mgwOp := sacloud.NewMobileGatewayOp(client)

	searched, err := mgwOp.Find(ctx, zone, &sacloud.FindCondition{})
	if err != nil {
		return false, fmt.Errorf("finding MobileGateway is failed: %s", err)
	}

	for _, mgw := range searched.MobileGateways {
		for _, iface := range mgw.Interfaces {
			if iface.SwitchScope != types.Scopes.Shared && iface.SwitchID == id {
				return true, nil
			}
		}
	}
	return false, nil
}
