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
	"os"
	"path/filepath"

	"github.com/sacloud/ftps"
)

func uploadFileViaFTPS(ctx context.Context, user, pass, host, file string) error {
	f, err := os.Open(file)
	if err != nil {
		return fmt.Errorf("opening file[%s] failed: %s", file, err)
	}
	defer f.Close() // nolint

	compCh := make(chan struct{})
	errCh := make(chan error)

	ftpClient := ftps.NewClient(user, pass, host)
	go func() {
		defer close(compCh)
		defer close(errCh)

		if err := ftpClient.UploadFile(filepath.Base(file), f); err != nil {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		f.Close() // nolint
		return ctx.Err()
	case err := <-errCh:
		return err
	case <-compCh:
		return nil
	}
}
