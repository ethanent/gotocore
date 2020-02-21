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
])

const buf = sc.build({
    "uname": -56,
    "tsts": 481324,
    "tbuf": Buffer.from([56, 69, 69, 69, 42, 0]),
})

console.log(Array.from(buf))
