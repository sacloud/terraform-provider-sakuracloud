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
