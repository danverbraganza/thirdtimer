thirdtimer.js: thirdtimer.go timefuncs/* components/*
	gopherjs build thirdtimer

local: thirdtimer.js
	gopherjs serve

.PHONY: local
