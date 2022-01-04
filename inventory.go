package vessels

type Inventory interface {
	Create() Vessel
	Get(...string) Vessel
}

// Initializer is an optional, life-cycle interface
// that an Inventory can implement to be initialized
// on register.
type Initializer interface {
	Initialize(*Universe)
}

func (u Universe) initialize(i Inventory) {
	if i, ok := i.(Initializer); ok {
		i.Initialize(&u)
	}
}

// Querier is an optional, life-cycle interface
// that an Inventory can implement to be queriable.
type Querier interface {
	Query(map[string][]string) ([]Vessel, error)
}

func query(i Inventory, s map[string][]string) (res interface{}) {
	var err error
	if q, ok := i.(Querier); ok {
		res, err = q.Query(s)
	} else {
		res, err = i, nil
	}
	if err != nil {
		panic(err)
	}
	return
}
