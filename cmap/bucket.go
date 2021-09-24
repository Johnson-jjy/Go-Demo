package cmap

type Bucket interface {
	Put(p pair, lock)
}
