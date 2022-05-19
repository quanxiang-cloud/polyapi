package timestamp

import (
	"fmt"
	"testing"
	"time"
)

func TestTimestamp(t *testing.T) {
	fmt.Println(time.Now().UTC().Format("2006-01-02T15:04:05Z"))
	fmt.Println(time.Now().Format("2006-01-02T15:04:05Z"))
	fmt.Println(supportedFormat)
	fmt.Println(ValidateTimeFormat(""))
	fmt.Println(ValidateTimeFormat("foo"))
	fmt.Println(Timestamp(""))
	fmt.Println(Timestamp("foo"))
	fmt.Println(Timestamp("YYYY-MM-DDThh:mm:ssZ"))
	fmt.Println(Timestamp("YYYY-MM-DDThh:mm:ss+0000"))
}
