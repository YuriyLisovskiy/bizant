// Copyright (c) 2018 Yuriy Lisovskiy
//
// Distributed under the BSD 3-Clause software license, see the accompanying
// file LICENSE or https://opensource.org/licenses/BSD-3-Clause.

// This LevelDB is implemented using https://github.com/syndtr/goleveldb
// source of Suryandaru Triandana <syndtr@gmail.com>

package leveldb

type Iterator interface {
	First() bool
	Last() bool
	Seek(key []byte) bool
	Next() bool
	Prev() bool
	Key() []byte
	Value() []byte
	Error() error
}

type IteratorGetter interface {
	Get(value []byte) (Iterator, error)
}

type EmptyIterator struct {
	Err error
}

func (*EmptyIterator) First() bool          { return false }
func (*EmptyIterator) Last() bool           { return false }
func (*EmptyIterator) Seek(key []byte) bool { return false }
func (*EmptyIterator) Next() bool           { return false }
func (*EmptyIterator) Prev() bool           { return false }
func (*EmptyIterator) Key() []byte          { return nil }
func (*EmptyIterator) Value() []byte        { return nil }
func (i *EmptyIterator) Error() error       { return i.Err }

type TwoLevelIterator struct {
	getter IteratorGetter
	index  Iterator
	data   Iterator
	err    error
}

func NewTwoLevelIterator(index Iterator, getter IteratorGetter) *TwoLevelIterator {
	return &TwoLevelIterator{getter: getter, index: index}
}

func (i *TwoLevelIterator) First() bool {
	if i.err != nil {
		return false
	}
	if !i.index.First() || !i.setData() {
		i.data = nil
		return false
	}
	return i.Next()
}

func (i *TwoLevelIterator) Last() bool {
	if i.err != nil {
		return false
	}
	if !i.index.Last() || !i.setData() {
		i.data = nil
		return false
	}
	if !i.data.Last() {
		i.data = nil
		return i.Prev()
	}
	return true
}

func (i *TwoLevelIterator) Seek(key []byte) bool {
	if i.err != nil {
		return false
	}
	if !i.index.Seek(key) || !i.setData() {
		i.data = nil
		return false
	}
	if !i.data.Seek(key) {
		return i.Next()
	}
	return true
}

func (i *TwoLevelIterator) Next() bool {
	if i.err != nil {
		return false
	}
	if i.data == nil || !i.data.Next() {
		if !i.index.Next() || !i.setData() {
			i.data = nil
			return false
		}
		return i.Next()
	}
	return true
}

func (i *TwoLevelIterator) Prev() bool {
	if i.err != nil {
		return false
	}
	if i.data == nil || !i.data.Prev() {
		if !i.index.Prev() || !i.setData() {
			i.data = nil
			return false
		}
		if !i.data.Last() {
			i.data = nil
			return i.Prev()
		}
		return true
	}
	return true
}

func (i *TwoLevelIterator) Key() []byte {
	if i.data == nil {
		return nil
	}
	return i.data.Key()
}
func (i *TwoLevelIterator) Value() []byte {
	if i.data == nil {
		return nil
	}
	return i.data.Value()
}
func (i *TwoLevelIterator) Error() error {
	if i.err != nil {
		return i.err
	} else if i.index.Error() != nil {
		return i.index.Error()
	} else if i.data != nil && i.data.Error() != nil {
		return i.data.Error()
	}
	return nil
}

func (i *TwoLevelIterator) setData() bool {
	i.data, i.err = i.getter.Get(i.index.Value())
	return i.err == nil
}
