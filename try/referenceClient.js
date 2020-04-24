var {Schema, StreamingAbstractor, types} = require("protocore")
const net = require("net")

const sc = new Schema([
	{
		"name": "tvarint",
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
	{
		"name": "tuint",
		"type": types.uint,
		"size": 16,
	},
])

const abs = new StreamingAbstractor()

abs.register('test', sc)

const conn = net.connect(8080)

conn.on('data', (b) => {
	console.log(b)
})

conn.on('connect', () => {
	console.log('ready')

	abs.bind(conn)

	abs.on('test', (data) => {
		console.log("Got test message.", data)
	})

	setInterval(() => {
		abs.send('test', {
			"tvarint": -535234,
			"tbuf": Buffer.from([51, 35, 51, 35, 64]),
			"tstr": "hey there :)",
			"tuint": 53,
		})

		console.log('sent')
	}, 4000)
})
