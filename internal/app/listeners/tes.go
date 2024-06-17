package listeners

import (
	"context"
	"fmt"
	"github.com/Marcellinom/tenant-management-saas/provider/event"
)

type TesListener struct{}

func NewTesListener() *TesListener {
	return &TesListener{}
}

func (t TesListener) Handle(ctx context.Context, payload event.Event) error {
	jsondata, _ := payload.JSON()
	fmt.Println(string(jsondata))
	return nil
}

func (t TesListener) MaxRetries() int {
	return 1
}

func (t TesListener) Name() string {
	return fmt.Sprintf("%T", t)
}
