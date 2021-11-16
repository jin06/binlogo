package node

import "github.com/jin06/binlogo/pkg/store/dao/dao_node"

func GetAllCount() (count int, err error){
	res, err := dao_node.AllNodes()
	if err != nil {
		return
	}
	count = len(res)
	return
}

func GetAllRegCount() (count int, err error) {
	res, err := dao_node.ALLRegisterNodes()
	if err != nil {
		return
	}
	count = len(res)
	return
}

