package gm

import "sync"

type Map[K comparable, V any] struct {
	m sync.Map
}

func (*Map[K, V]) cast(v any, ok bool) (V, bool) {
	if ok {
		return v.(V), true
	}
	var z V
	return z, false
}

func (g *Map[K, V]) Load(k K) (V, bool) {
	return g.cast(g.m.Load(k))
}

func (g *Map[K, V]) Store(k K, v V) {
	g.m.Store(k, v)
}

func (g *Map[K, V]) Clear() {
	g.m.Clear()
}

func (g *Map[K, V]) LoadOrStore(k K, v V) (V, bool) {
	return g.cast(g.m.LoadOrStore(k, v))
}

func (g *Map[K, V]) LoadAndDelete(k K) (V, bool) {
	return g.cast(g.m.LoadAndDelete(k))
}

func (g *Map[K, V]) Delete(k K) {
	g.m.Delete(k)
}

func (g *Map[K, V]) Swap(k K, v V) (V, bool) {
	return g.cast(g.m.Swap(k, v))
}

func (g *Map[K, V]) CompareAndSwap(k K, old V, new V) bool {
	return g.m.CompareAndSwap(k, old, new)
}

func (g *Map[K, V]) CompareAndDelete(k K, v V) bool {
	return g.m.CompareAndDelete(k, v)
}

func (g *Map[K, V]) Range(f func(k K, v V) bool) {
	g.m.Range(func(k, v any) bool {
		return f(k.(K), v.(V))
	})
}
