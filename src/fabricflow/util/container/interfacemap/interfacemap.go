// -*- coding: utf-8 -*-

package interfacemap

func SelectOrInsert(m map[interface{}]interface{}, names ...string) (map[interface{}]interface{}, bool) {
	var child interface{} = m
	for _, name := range names {
		parent, ok := child.(map[interface{}]interface{})
		if !ok {
			return nil, false
		}

		child, ok = parent[name]
		if !ok {
			child = map[interface{}]interface{}{}
			parent[name] = child
		}
	}

	m, ok := child.(map[interface{}]interface{})
	return m, ok
}

func Select(m map[interface{}]interface{}, names ...string) (interface{}, bool) {
	var child interface{} = m
	for _, name := range names {
		parent, ok := child.(map[interface{}]interface{})
		if !ok {
			return nil, false
		}

		child, ok = parent[name]
		if !ok {
			return nil, false
		}
	}

	return child, true
}

func SelectList(m map[interface{}]interface{}, names ...string) ([]interface{}, bool) {
	node, ok := Select(m, names...)
	if !ok {
		return nil, false
	}

	nodeList, ok := node.([]interface{})
	return nodeList, ok
}

func SelectMap(m map[interface{}]interface{}, names ...string) (map[interface{}]interface{}, bool) {
	node, ok := Select(m, names...)
	if !ok {
		return nil, false
	}

	nodeMap, ok := node.(map[interface{}]interface{})
	return nodeMap, ok
}

func Set(m map[interface{}]interface{}, value interface{}, names ...string) bool {
	index := len(names) - 1
	if index < 0 {
		return false
	}

	parentNode, ok := SelectOrInsert(m, names[:index]...)
	if !ok {
		return false
	}

	parentNode[names[index]] = value
	return true
}

func Remove(m map[interface{}]interface{}, names ...string) bool {
	index := len(names) - 1
	if index < 0 {
		return false
	}

	parentNode, ok := SelectMap(m, names[:index]...)
	if !ok {
		return false
	}

	delete(parentNode, names[index])
	return true
}
