module github.com/lekuruu/hexagon

go 1.22.7

require (
	github.com/lekuruu/hexagon/common v0.0.0-20241005191625-f1f1499b1fd1
	github.com/lekuruu/hexagon/hnet v0.0.0-00010101000000-000000000000
)

require (
	github.com/soniakeys/meeus/v3 v3.0.1 // indirect
	github.com/soniakeys/unit v1.0.0 // indirect
)

require (
	github.com/lekuruu/go-raknet v0.0.0-20241003121121-43b332df40f2 // indirect
	github.com/lekuruu/hexagon/hscore v0.0.0-20241004135010-1d8d562aeb0e
	golang.org/x/exp v0.0.0-20240909161429-701f63a606c0 // indirect
)

replace github.com/lekuruu/hexagon/hnet => ./hnet

replace github.com/lekuruu/hexagon/common => ./common

replace github.com/lekuruu/hexagon/hscore => ./hscore
