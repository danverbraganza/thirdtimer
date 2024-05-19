thirdtimer.js: thirdtimer.go
	gopherjs build thirdtimer

local: thirdtimer.js
	gopherjs serve

.PHONY: local
