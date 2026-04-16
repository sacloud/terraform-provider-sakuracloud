// Copyright 2016-2025 terraform-provider-sakuracloud authors
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

package ftps

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/jlaffaye/ftp"
)

func UploadFile(ctx context.Context, user, pass, host, file string) error {
	f, err := os.Open(filepath.Clean(file))
	if err != nil {
		return fmt.Errorf("opening file[%s] failed: %s", file, err)
	}
	defer f.Close() //nolint

	compCh := make(chan struct{})
	errCh := make(chan error)

	go func() {
		defer close(compCh)
		defer close(errCh)

		log.Printf("[INFO] upload file to ftps %s", host)
		conn, err := ftp.Dial(
			fmt.Sprintf("%s:%d", host, 21),
			ftp.DialWithTimeout(30*time.Minute),
			ftp.DialWithExplicitTLS(&tls.Config{
				ServerName: host,
				MinVersion: tls.VersionTLS12,
				MaxVersion: tls.VersionTLS13,
			}))
		if err != nil {
			errCh <- fmt.Errorf("failed to connect to FTP server[%s]: %w", host, err)
			return
		}
		defer conn.Quit() //nolint:errcheck

		if err := conn.Login(user, pass); err != nil {
			errCh <- fmt.Errorf("failed to login to FTP server[%s]: %w", host, err)
			return
		}

		if err := conn.Stor(filepath.Base(file), f); err != nil {
			errCh <- fmt.Errorf("failed to upload file[%s]: %w", host, err)
			return
		}

		compCh <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		f.Close() //nolint
		return ctx.Err()
	case err := <-errCh:
		return err
	case <-compCh:
		return nil
	}
}
