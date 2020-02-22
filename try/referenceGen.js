var protocore = require("protocore")

const sc = new protocore.Schema([
	{
		"name": "uname",
		"type": protocore.types.varint,
	},
	{
		"name": "tsts",
		"type": protocore.types.varint,
	},
	{
		"name": "tbuf",
		"type": protocore.types.buffer,
	},
	{
		"name": "tstr",
		"type": protocore.types.string,
	},
	{
		"name": "tuint",
		"type": protocore.types.uint,
		"size": 16,
	},
])

const buf = sc.build({
	"uname": -56,
	"tsts": 481324,
	"tbuf": Buffer.from([56, 69, 69, 69, 42, 0]),
	"tstr": "HeY THEre! 3546",
	"tuint": 533,
})

console.log(Array.from(buf))
