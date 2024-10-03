module github.com/lekuruu/hexagon

go 1.22.7

require (
	github.com/lekuruu/hexagon/common v0.0.0-20241003124433-ad9d5dc7b056
	github.com/lekuruu/hexagon/hnet v0.0.0-00010101000000-000000000000
)

require github.com/lekuruu/go-raknet v0.0.0-20241003121121-43b332df40f2 // indirect

replace github.com/lekuruu/hexagon/hnet => ./hnet

replace github.com/lekuruu/hexagon/common => ./common
