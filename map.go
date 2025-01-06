package event

import "sync"

type gm[K comparable, V any] struct {
	m sync.Map
}

func (*gm[K, V]) cast(v any, ok bool) (V, bool) {
	if ok {
		return v.(V), true
	}
	var z V
	return z, false
}

func (g *gm[K, V]) Load(k K) (V, bool) {
	return g.cast(g.m.Load(k))
}

func (g *gm[K, V]) Store(k K, v V) {
	g.m.Store(k, v)
}

func (g *gm[K, V]) Clear() {
	g.m.Clear()
}

func (g *gm[K, V]) LoadOrStore(k K, v V) (V, bool) {
	return g.cast(g.m.LoadOrStore(k, v))
}

func (g *gm[K, V]) LoadAndDelete(k K) (V, bool) {
	return g.cast(g.m.LoadAndDelete(k))
}

func (g *gm[K, V]) Delete(k K) {
	g.m.Delete(k)
}

func (g *gm[K, V]) Swap(k K, v V) (V, bool) {
	return g.cast(g.m.Swap(k, v))
}

func (g *gm[K, V]) CompareAndSwap(k K, old V, new V) bool {
	return g.m.CompareAndSwap(k, old, new)
}

func (g *gm[K, V]) CompareAndDelete(k K, v V) bool {
	return g.m.CompareAndDelete(k, v)
}

func (g *gm[K, V]) Range(f func(k K, v V) bool) {
	g.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}
