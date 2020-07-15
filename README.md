# CBOR Utility

This utility can decode data in a CBOR payload located at a specific data path.

This utility can also be used to update data in a CBOR payload at a specific data path.

---

## Building

```
make build
```

This builds a binary in `bin/cborutil`

## Usage

### Decoding CBOR Data

```
bin/cborutil decode -cbor <payload> [ -path data.path -base64 -json ]
```

- `-cbor` CBOR payload in hex or base64 format
- `-path` Data path to decode using dot notation
- `-base64` Payload is in base64 format
- `-json` Return decoded result in JSON format?

#### Examples

##### Decoding from hex:

```
$ bin/cborutil decode -cbor a26161016162820203 -path b.0
2
```

##### Decoding from base64:

```
$ bin/cborutil decode -cbor omFhAWFiggID -base64 -path b.0
2
```

##### Decoding to JSON:

```
$ bin/cborutil decode -cbor omFhAWFiggID -base64 -path b -json
[2,3]
```

Can you also omit the `-path` arg to decode the whole payload:

```
$ bin/cborutil decode -cbor omFhAWFiggID -base64 -json
{"a":1,"b":[2,3]}
```

### Updating CBOR Data

```
bin/cborutil update -cbor <payload> -data <input> -path data.path [ -base64 -json ]
```

- `-cbor` CBOR payload in hex or base64 format
- `-data` Input data
- `-path` Data path to update using dot notation
- `-base64` Payload is in base64 format
- `-json` Input data is in JSON format?

#### Examples

##### Updating a hex payload

```
$ bin/cborutil update -cbor a26161016162820203 -data "changed me" -path a
a261616a6368616e676564206d656162820203
```

##### Updating a base64 payload

```
$ bin/cborutil update -cbor omFhAWFiggID -base64 -data "changed me" -path a
omFhamNoYW5nZWQgbWVhYoICAw==
```

##### Updating using JSON

```
$ bin/cborutil update -cbor omFhAWFiggID -base64 -data '{"aa":"AAA"}' -json -path a
omFhoWJhYWNBQUFhYoICAw==
```

## Dot Notation

This utility uses dot notation to specify the path to a node in the data. For example, let's take a look at the following JSON data:

```
{
  "Student": {
    "Name": "Eve",
    "Parents": [
      {
        "Name": "Alice"
      },
      {
        "Name": "Bob"
      }
    ]
  }
}
```

Here, if we want to get the student's name, we could reference `Student.Name`

Arrays are referenced by their integer index, starting from 0. For example to get the Student's Parent's names, we would use `Student.Parents.0.Name` and `Student.Parents.1.Name`

Let's use the utility to update the CBOR data for Eve. Her father wants his name to be listed as Robert instead of Bob.

```
$ bin/cborutil update \
-cbor oWdTdHVkZW50omROYW1lY0V2ZWdQYXJlbnRzgqFkTmFtZWVBbGljZaFkTmFtZWNCb2I= \
-base64 -path Student.Parents.1.Name -data Robert
oWdTdHVkZW50omROYW1lY0V2ZWdQYXJlbnRzgqFkTmFtZWVBbGljZaFkTmFtZWZSb2JlcnQ=
```

We can also use the utility to confirm the update happened correctly:

```
$ bin/cborutil decode \
-cbor oWdTdHVkZW50omROYW1lY0V2ZWdQYXJlbnRzgqFkTmFtZWVBbGljZaFkTmFtZWZSb2JlcnQ= \
-base64 -json
{"Student":{"Name":"Eve","Parents":[{"Name":"Alice"},{"Name":"Robert"}]}}
```

**Note** You can couple `cborutil update` with the `-json` flag to not just update individual values, but to replace entire portions of the data tree.

For example, we entered the wrong student data. Let's go ahead and correct it:

```
$ bin/cborutil update \
-cbor oWdTdHVkZW50omROYW1lY0V2ZWdQYXJlbnRzgqFkTmFtZWVBbGljZaFkTmFtZWNCb2I= \
-base64 -path Student \
-data '{"Name":"John","Interests":["Gaming","Blockchain","Piano"]}' \
-json
oWdTdHVkZW50omROYW1lZEpvaG5pSW50ZXJlc3Rzg2ZHYW1pbmdqQmxvY2tjaGFpbmVQaWFubw==
```

And again we can use the utility to verify the data tree was updated:

```
$ bin/cborutil decode \
-cbor oWdTdHVkZW50omROYW1lZEpvaG5pSW50ZXJlc3Rzg2ZHYW1pbmdqQmxvY2tjaGFpbmVQaWFubw== \
-base64 -json
{"Student":{"Interests":["Gaming","Blockchain","Piano"],"Name":"John"}}
```
