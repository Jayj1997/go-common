/*
 * @Author       : jayj
 * @Date         : 2021-09-13 10:02:48
 * @Description  :
 */
package btree

import (
	"encoding/json"

	"github.com/Jayj1997/go-common"
)

// ToJSON outputs the JSON representation of the tree
func (tree *Tree) ToJSON() ([]byte, error) {

	elements := make(map[string]interface{})
	it := tree.Iterator()

	for it.Next() {
		elements[common.ToString(it.Key())] = it.Value()
	}

	return json.Marshal(&elements)
}

// FromJSON populates the tree from the input JSON representation
func (tree *Tree) FromJSON(data []byte) error {

	elements := make(map[string]interface{})
	err := json.Unmarshal(data, &elements)

	if err == nil {

		tree.Clear()

		for key, value := range elements {

			tree.Insert(key, value)
		}
	}

	return err
}
