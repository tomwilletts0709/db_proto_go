import (
	"os"
)

type Aof struct {
	file *os.File
	rd *bufio.Reader
	my sync.Mutex
}

func New Aof(path string) (*Aof, error) {
	f, err := os.OpenFile(path, os.O_)
}