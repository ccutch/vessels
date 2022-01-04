package memory

import (
	"fmt"
	"strings"
	"vessels"
)

type Memory map[string]*Buffer

func Empty() Memory {
	return Memory(make(map[string]*Buffer))
}

func (m Memory) Create() vessels.Vessel {
	var n = fmt.Sprintf("v%d", len(m)+1)
	m[n] = new(Buffer)
	return m[n]
}

func (m Memory) Get(s ...string) vessels.Vessel {
	var n = strings.Join(s, ":")
	if _, ok := m[n]; !ok {
		m[n] = new(Buffer)
	}
	return m[n]
}
