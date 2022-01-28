package sdk

//// TODO (ahawker) I dont think this can be changed atm in the wasmer.
//const namespace = "env"
//
//// Exports returns the namespace and functions that should be exported
//// as part of the SDK.
//func Exports() (string, []Function) {
//	return namespace, host.Exports2
//	//return namespace, exportsToFunctions(host.Exports)
//}

//// exportsToFunctions creates a slice of sdk.Function instances for the
//// given map of host.
//func exportsToFunctions(exported map[string]host.HostFunc) []Function {
//	var functions []Function
//
//	for name, export := range exported {
//		functions = append(functions, NewFunction(name, export))
//	}
//	return functions
//}
