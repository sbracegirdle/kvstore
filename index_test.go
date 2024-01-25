package main

import (
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	tree := createIndexTree([]IndexValue{
		{1, 99},
		{2, 88},
		{3, 77},
		{4, 66},
		{5, 55},
	}, 3)

	tests := []struct {
		key      uint32
		expected bool
	}{
		{1, true},
		{2, true},
		{6, false},
		{0, false},
		{4, true},
	}

	for _, test := range tests {
		result, _ := tree.Root.Search(test.key)
		if (result != nil) != test.expected {
			t.Errorf("Expected %v, got %v", test.expected, result != nil)
		}
	}
}

// Assuming a GetIndexValues function exists
func getIndexValues(node *IndexTreeNode, keys *[]IndexValue) {
	if node != nil {
		for _, child := range node.Child {
			getIndexValues(child, keys)
		}

		*keys = append(*keys, node.Keys...)
	}
}

func TestInsert(t *testing.T) {
	tests := []struct {
		keys     []IndexValue
		expected []IndexValue
	}{
		{[]IndexValue{
			{1, 99},
			{2, 88},
			{3, 77},
			{4, 66},
			{5, 55},
		}, []IndexValue{
			{1, 99},
			{2, 88},
			{3, 77},
			{4, 66},
			{5, 55},
		}},
		{[]IndexValue{
			{5, 99},
			{4, 88},
			{3, 77},
			{2, 66},
			{1, 55},
		}, []IndexValue{
			{1, 55},
			{2, 66},
			{3, 77},
			{4, 88},
			{5, 99},
		}},
		{[]IndexValue{}, []IndexValue{}},
		{[]IndexValue{{1, 1}}, []IndexValue{{1, 1}}},
	}

	for _, test := range tests {
		tree := &IndexTree{
			Root: &IndexTreeNode{
				IsLeaf: true,
				Keys:   make([]IndexValue, 0),
				Child:  make([]*IndexTreeNode, 0),
			},
			MinDegree: 3,
		}

		for _, key := range test.keys {
			tree.Insert(&key)
		}

		var keys []IndexValue
		getIndexValues(tree.Root, &keys)

		if test.expected == nil || len(test.expected) == 0 {
			if len(keys) > 0 {
				t.Errorf("Expected %v, got %v", test.expected, keys)
			}
		} else if !reflect.DeepEqual(keys, test.expected) {
			t.Errorf("Expected %v, got %v", test.expected, keys)
		}
	}
}
