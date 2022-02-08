package acq

import "fmt"

const help = ` 
-*-=============================== pre_check version 1.0 ===============================-*
Usage of ./pre_check:
	-ClusterName string
		input cluster name(mycluster). default mycluster (default "none")
    -OpType string 
		operation type(precheck, switchover). default precheck (default "none")
Example:
	./pre_check --OpType=precheck --ClusterName=mycluster

`

func ByDefaultAction() {
	fmt.Println(help)
}
