// Copyright 2016-2021 The Libsacloud Authors
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

package power

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/sacloud/libsacloud/v2/sacloud"
	"github.com/sacloud/libsacloud/v2/sacloud/accessor"
	"github.com/sacloud/libsacloud/v2/sacloud/types"
)

const (
	// DefaultBootRetrySpan BootRetrySpanのデフォルト値
	DefaultBootRetrySpan = 20 * time.Second
	// DefaultShutdownRetrySpan ShutdownRetrySpanのデフォルト値
	DefaultShutdownRetrySpan = 20 * time.Second
)

var (
	// BootRetrySpan 起動APIをコールしてからリトライするまでの待機時間
	BootRetrySpan = DefaultBootRetrySpan

	// ShutdownRetrySpan シャットダウンAPIをコールしてからリトライするまでの待機時間
	ShutdownRetrySpan = DefaultShutdownRetrySpan
)

/************************************************
 * Server
 ***********************************************/

// BootServer 起動
func BootServer(ctx context.Context, client ServerAPI, zone string, id types.ID) error {
	return boot(ctx, &serverHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	})
}

// ShutdownServer シャットダウン
func ShutdownServer(ctx context.Context, client ServerAPI, zone string, id types.ID, force bool) error {
	return shutdown(ctx, &serverHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	}, force)
}

/************************************************
 * LoadBalancer
 ***********************************************/

// BootLoadBalancer 起動
func BootLoadBalancer(ctx context.Context, client LoadBalancerAPI, zone string, id types.ID) error {
	return boot(ctx, &loadBalancerHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	})
}

// ShutdownLoadBalancer シャットダウン
func ShutdownLoadBalancer(ctx context.Context, client LoadBalancerAPI, zone string, id types.ID, force bool) error {
	return shutdown(ctx, &loadBalancerHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	}, force)
}

/************************************************
 * Database
 ***********************************************/

// BootDatabase 起動
func BootDatabase(ctx context.Context, client DatabaseAPI, zone string, id types.ID) error {
	return boot(ctx, &databaseHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	})
}

// ShutdownDatabase シャットダウン
func ShutdownDatabase(ctx context.Context, client DatabaseAPI, zone string, id types.ID, force bool) error {
	return shutdown(ctx, &databaseHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	}, force)
}

/************************************************
 * VPCRouter
 ***********************************************/

// BootVPCRouter 起動
func BootVPCRouter(ctx context.Context, client VPCRouterAPI, zone string, id types.ID) error {
	return boot(ctx, &vpcRouterHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	})
}

// ShutdownVPCRouter シャットダウン
func ShutdownVPCRouter(ctx context.Context, client VPCRouterAPI, zone string, id types.ID, force bool) error {
	return shutdown(ctx, &vpcRouterHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	}, force)
}

/************************************************
 * NFS
 ***********************************************/

// BootNFS 起動
func BootNFS(ctx context.Context, client NFSAPI, zone string, id types.ID) error {
	return boot(ctx, &nfsHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	})
}

// ShutdownNFS シャットダウン
func ShutdownNFS(ctx context.Context, client NFSAPI, zone string, id types.ID, force bool) error {
	return shutdown(ctx, &nfsHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	}, force)
}

/************************************************
 * MobileGateway
 ***********************************************/

// BootMobileGateway 起動
func BootMobileGateway(ctx context.Context, client MobileGatewayAPI, zone string, id types.ID) error {
	return boot(ctx, &mobileGatewayHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	})
}

// ShutdownMobileGateway シャットダウン
func ShutdownMobileGateway(ctx context.Context, client MobileGatewayAPI, zone string, id types.ID, force bool) error {
	return shutdown(ctx, &mobileGatewayHandler{
		ctx:    ctx,
		client: client,
		zone:   zone,
		id:     id,
	}, force)
}

type handler interface {
	boot() error
	shutdown(force bool) error
	read() (interface{}, error)
}

func boot(ctx context.Context, h handler) error {
	if err := h.boot(); err != nil {
		return err
	}

	retryTimer := time.NewTicker(BootRetrySpan)
	defer retryTimer.Stop()

	inProcess := false

	waiter := sacloud.WaiterForUp(h.read)
	compCh, progressCh, errCh := waiter.AsyncWaitForState(ctx)

	var state interface{}

	for {
		select {
		case <-ctx.Done():
			return errors.New("canceled")
		case <-compCh:
			return nil
		case s := <-progressCh:
			state = s
		case <-retryTimer.C:
			if inProcess {
				continue
			}
			if state != nil && state.(accessor.InstanceStatus).GetInstanceStatus().IsDown() {
				if err := h.boot(); err != nil {
					if err, ok := err.(sacloud.APIError); ok {
						if err.ResponseCode() == http.StatusConflict {
							inProcess = true
							continue
						}
					}
					return err
				}
			}
		case err := <-errCh:
			return err
		}
	}
}

func shutdown(ctx context.Context, h handler, force bool) error {
	if err := h.shutdown(force); err != nil {
		return err
	}

	retryTimer := time.NewTicker(ShutdownRetrySpan)
	defer retryTimer.Stop()

	inProcess := false

	waiter := sacloud.WaiterForDown(h.read)
	compCh, progressCh, errCh := waiter.AsyncWaitForState(ctx)

	var state interface{}

	for {
		select {
		case <-compCh:
			return nil
		case s := <-progressCh:
			state = s
		case <-retryTimer.C:
			if inProcess {
				continue
			}
			if state != nil && state.(accessor.InstanceStatus).GetInstanceStatus().IsUp() {
				if err := h.shutdown(force); err != nil {
					if err, ok := err.(sacloud.APIError); ok {
						if err.ResponseCode() == http.StatusConflict {
							inProcess = true
							continue
						}
					}
					return err
				}
			}
		case err := <-errCh:
			return err
		}
	}
}
