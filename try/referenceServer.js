var {Schema, StreamingAbstractor, types} = require("protocore")

const sc = new Schema([
	{
		"name": "uname",
		"type": types.varint,
	},
	{
		"name": "tsts",
		"type": types.varint,
	},
	{
		"name": "tbuf",
		"type": types.buffer,
	},
	{
		"name": "tstr",
		"type": types.string,
	},
])

const abs = new StreamingAbstractor()

abs.register('tester', sc)

const buf = sc.build({
	"uname": -56,
	"tsts": 481324,
	"tbuf": Buffer.from([56, 69, 69, 69, 42, 0]),
	"tstr": "HeY THEre! 3546",
})

console.log(Array.from(buf))
