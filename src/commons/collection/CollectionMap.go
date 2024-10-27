package collection

type CollectionMap[T comparable, K any] struct {
	items map[T]K
}

func FromMap[T comparable, K any](items map[T]K) *CollectionMap[T, K] {
	return &CollectionMap[T, K]{
		items,
	}
}

func EmptyMap[T comparable, K any]() *CollectionMap[T, K] {
	return FromMap(make(map[T]K))
}

func (c *CollectionMap[T, K]) Size() int {
	return len(c.items)
}

func (c *CollectionMap[T, K]) Find(key T) (*K, bool) {
	value, exists := c.items[key]
	return &value, exists
}

func (c *CollectionMap[T, K]) FindOnePredicate(predicate func(T, K) bool) (*K, bool) {
	items := c.FindPredicate(predicate)
	if len(items) == 0 {
		return nil, false
	}
	return &items[0], true
}

func (c *CollectionMap[T, K]) FindPredicate(predicate func(T, K) bool) []K {
	filter := []K{}
	for k, v := range c.items {
		if predicate(k, v) {
			filter = append(filter, v)
		}
	}
	return filter
}

func (c *CollectionMap[T, K]) Filter(predicate func(T, K) bool) *CollectionMap[T, K] {
	filter := map[T]K{}
	for key, item := range c.items {
		if predicate(key, item) {
			filter[key] = item
		}
	}
	c.items = filter
	return c
}

func (c *CollectionMap[T, K]) MapSelf(predicate func(T, K) K) *CollectionMap[T, K] {
	return MapMap(c, predicate)
}

func (c *CollectionMap[T, K]) Exists(key T) bool {
	_, exists := c.items[key]
	return exists
}

func (c *CollectionMap[T, K]) Put(key T, item K) (*K, bool) {
	old, exists := c.Find(key)
	c.items[key] = item
	return old, exists
}

func (c *CollectionMap[T, K]) PutIfAbsent(key T, item K) (*K, bool) {
	old, exists := c.Find(key)
	if !exists {
		c.Put(key, item)
	}
	return old, exists
}

func (c *CollectionMap[T, K]) Merge(other map[T]K) *CollectionMap[T, K] {
	for key := range other {
		c.items[key] = other[key]
	}
	return c
}

func (c *CollectionMap[T, K]) Remove(key T, item K) (*K, bool) {
	old, exists := c.Find(key)
	delete(c.items, key)
	return old, exists
}

func (collection CollectionMap[T, K]) Keys() []T {
	keys := make([]T, 0, len(collection.items))
	for key := range collection.items {
		keys = append(keys, key)
	}
	return keys
}

func (collection CollectionMap[T, K]) KeysCollection() *CollectionList[T] {
	return FromList(collection.Keys())
}

func (c *CollectionMap[T, K]) Values() []K {
	values := make([]K, 0, len(c.items))
	for key := range c.items {
		values = append(values, c.items[key])
	}
	return values
}

func (collection CollectionMap[T, K]) ValuesCollection() *CollectionList[K] {
	return FromList(collection.Values())
}

func (c *CollectionMap[T, K]) ValuesInterface() []any {
	values := make([]any, 0, len(c.items))
	for key := range c.items {
		values = append(values, c.items[key])
	}
	return values
}

func (c *CollectionMap[T, K]) Pairs() []Pair[T, K] {
	pairs := make([]Pair[T, K], 0, len(c.items))
	for key := range c.items {
		pairs = append(pairs, Pair[T, K]{
			key:   key,
			value: c.items[key],
		})
	}
	return pairs
}

func (c *CollectionMap[T, K]) Clean() *CollectionMap[T, K] {
	c.items = make(map[T]K)
	return c
}

func (c *CollectionMap[T, K]) Clone() *CollectionMap[T, K] {
	return FromMap(c.items)
}

func (collection CollectionMap[T, K]) Collect() map[T]K {
	return collection.items
}

func ListToMapFrom[T comparable, K, E any](list []K, predicate func(K) (T, E)) *CollectionMap[T, E] {
	mapped := map[T]E{}
	for _, item := range list {
		t, e := predicate(item)
		mapped[t] = e
	}
	return &CollectionMap[T, E]{
		items: mapped,
	}
}

func MapMapFrom[T comparable, K, E any](items map[T]K, predicate func(T, K) E) *CollectionMap[T, E] {
	return MapMap(FromMap(items), predicate)
}

func MapMap[T comparable, K, E any](c *CollectionMap[T, K], predicate func(T, K) E) *CollectionMap[T, E] {
	mapped := map[T]E{}
	for key, item := range c.items {
		mapped[key] = predicate(key, item)
	}
	return &CollectionMap[T, E]{
		items: mapped,
	}
}

func MapMerge[T comparable, K any](target, source map[T]K) map[T]K {
	for k, v := range source {
		target[k] = v
	}
	return target
}
