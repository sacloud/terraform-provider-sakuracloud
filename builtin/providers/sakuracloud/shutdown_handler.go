package sakuracloud

import (
	"github.com/hashicorp/terraform/helper/schema"
	"time"
)

type shutdownHandler interface {
	Stop(id int64) (bool, error)
	Shutdown(id int64) (bool, error)
	SleepUntilDown(id int64, timeout time.Duration) error
}

var (
	powerManageTimeoutKey   = "power_manage_timeout"
	powerManageTimeoutParam = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  60,
	}
	powerManageTimeoutParamForceNew = &schema.Schema{
		Type:     schema.TypeInt,
		Optional: true,
		Default:  60,
		ForceNew: true,
	}
)

func handleShutdown(handler shutdownHandler, id int64, d *schema.ResourceData, defaultTimeOut time.Duration) error {

	timeout := defaultTimeOut
	if v, ok := d.GetOk(powerManageTimeoutKey); ok {
		s := v.(int)
		timeout = time.Duration(s) * time.Second
	}

	// graceful shutdown
	_, err := handler.Shutdown(id)
	if err != nil {
		return err
	}

	// wait
	if err = handler.SleepUntilDown(id, timeout); err != nil {
		// force shutdown
		if _, err := handler.Stop(id); err != nil {
			return err
		}
	}

	return handler.SleepUntilDown(id, timeout)
}
