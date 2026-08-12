package cellwork

import "runtime"

type runtimeMem struct{ total uint64 }

func readMem(m *runtimeMem) {
	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	m.total = ms.TotalAlloc
}
