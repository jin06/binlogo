package node

import "github.com/jin06/binlogo/v2/pkg/store/dao/dao_node"

// GetAllCount get all nodes count
func GetAllCount() (count int, err error) {
	res, err := dao_node.AllNodes()
	if err != nil {
		return
	}
	if res == nil {
		return 0, nil
	}
	count = len(res)
	return
}

// GetAllRegCount get all nodes count that already registered
func GetAllRegCount() (count int, err error) {
	res, err := dao_node.ALLRegisterNodes()
	if err != nil {
		return
	}
	count = len(res)
	return
}
