package limitrate_test

import (
	"fmt"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/quanxiang-cloud/polyapi/pkg/core/gate/limitrate"
)

func TestBucket(t *testing.T) {
	b := limitrate.NewBucket(37)
	fmt.Println(b.ShowTokenArg())
	b.Start()
	fmt.Println(b.ShowTokenArg())
	time.Sleep(time.Millisecond * 100)
	var wg sync.WaitGroup
	req := func(id int) {
		wg.Add(1)
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(id*10+60)+30))
		fmt.Println(id, b.ShowTokenArg(), b.GetToken())
		wg.Done()
	}
	for i := 1; i <= 10; i++ {
		go req(i)
	}
	wg.Wait()
	b.Stop()
}
