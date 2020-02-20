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
])

const buf = sc.build({
    "uname": 568,
    "tsts": 481324,
})

console.log(Array.from(buf))
