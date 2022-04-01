# mockaroo2hl7v2

This is a simple tool to generate HL7v2 messages from a specific [Mockaroo](https://mockaroo.com) JSON output format.  I needed a way to generate fake HL7v2 messages to test things like Mirth channels + wanted to play around with [Go](https://go.dev).

The output HL7v2 is generating using Go [Text Templates](https://pkg.go.dev/text/template).  You can specify a template file when running, if one is not specified there is one in the [main.go](main.go) file.

## The Input

This utility expects a very specific JSON input, one which I use Mockaroo to setup.  There is a backup of the Mockaroo _schema_ in this repo named [HL7JSON.mockaroo.json](HL7JSON.mockaroo.json).  

A simple example:

```json
{
  "Patient": {
    "Ids": [
      {
        "IdNumber": "PG23FK030030",
        "CheckDigit": "2",
        "CheckDigitScheme": "BCV",
        "IdType": "MR"
      },
      {
        "IdNumber": "TO64MN065772",
        "CheckDigit": "2",
        "CheckDigitScheme": "BCV",
        "IdType": "AN"
      }
    ],
    "Name": {
      "Family": "PatricksonZZFAKE",
      "Given": "ZZFAKELillis",
      "Second": "I",
      "Suffix": null,
      "Prefix": null
    },
    "MotherMaidenName": "Tallboy",
    "DOB": "1998-01-18T16:48:01Z",
    "Gender": "F",
    "Race": {
      "Id": "2076-8",
      "Text": "Native Hawaiian or Other Pacific Islander",
      "CodingSystem": "HL70005"
    },
    "Address": [
      {
        "Line1": "76255 Harbort Park",
        "Line2": "Number 18",
        "City": "Memphis",
        "State": "TN",
        "Zip": "38136",
        "Country": "US",
        "Type": "M"
      }
    ],
    "Comm": [
      {
        "UseCode": "PRN",
        "Type": "CP",
        "PhoneNumber": "5791303873",
        "Email": "",
        "CountryCode": "+234",
        "AreaCode": "579",
        "LocalNumber": "1303873",
        "UnformattedTelephoneNumber": "2345791303873"
      },
      {
        "UseCode": "WPN",
        "Type": "CP",
        "PhoneNumber": "5689724029",
        "Email": "",
        "CountryCode": "+62",
        "AreaCode": "568",
        "LocalNumber": "9724029",
        "UnformattedTelephoneNumber": "625689724029"
      }
    ],
    "Language": {
      "Id": "pl",
      "Text": "Polish",
      "CodingSystem": "urn:ietf:bcp:47"
    },
    "MaritalStatus": {
      "Id": "O",
      "Text": "Other",
      "CodingSystem": "HL70002"
    },
    "Religion": {
      "Id": "HOT",
      "Text": "Hindu: Other",
      "CodingSystem": "HL70006"
    },
    "SSN": "661-41-7685",
    "DriversLicense": {
      "Number": "1913688611",
      "State": "TN",
      "ExpirationDate": "2023-03-28T16:48:01Z"
    },
    "EthnicGroup": {
      "Id": "N",
      "Text": "Not Hispanic or Latino",
      "CodingSystem": "HL70189"
    },
    "MultipleBirthIndicator": "Y",
    "BirthOrder": 5,
    "DeathIndicator": "Y",
    "DeathDateTime": "2022-03-27T16:48:01Z"
  },
  "Providers": [
    {
      "Id": {
        "IdNumber": "714CW28081",
        "CheckDigit": "W",
        "CheckDigitScheme": "NPI",
        "IdType": "NPI"
      },
      "Degree": "HS",
      "Name": {
        "Family": "O' Bee",
        "Given": "Jacklin",
        "Second": "B",
        "Suffix": null,
        "Prefix": null
      },
      "NameTypeCode": "B",
      "Role": {
        "Id": "PP",
        "Text": "Primary Care Provider",
        "CodingSystem": "HL70443"
      }
    }
  ],
  "PatientClass": {
    "Id": "C",
    "Text": "Commercial Account",
    "CodingSystem": "HL70004"
  },
  "Event": "A07"
}
```

The input can be a single JSON object (as in the example above), or an array of such objects.  A single object input is going to produce a single HL7 message, while an array is going to produce as many HL7 messages as in the array.

See later in this document for a few pointers on importing the Mockaroo schema to your own Mockaroo account.

## The Output

Expected output if you run the example JSON above through mockaroo2hl7v2 would be:

```
MSH|^~\&|Mirth|Hospital|HIE|HIE|20220329212610||ADT^A07|afd629ec-0288-470e-8b1e-69d4458a00c9|P|2.5.1
EVN|A07|20220329212610
PID|1||PG23FK030030^2^BCV^^MR~TO64MN065772^2^BCV^^AN||PatricksonZZFAKE^ZZFAKELillis^I^^|Tallboy|19980118|F||2076-8^Native Hawaiian or Other Pacific Islander^HL70005|76255 Harbort Park^Number 18^Memphis^TN^38136^US^M||5791303873^PRN^CP^^+234^579^1303873^^^^^2345791303873|5689724029^WPN^CP^^+62^568^9724029^^^^^625689724029|pl|O^Other^HL70002|HOT|TO64MN065772^2^BCV^^AN|661-41-7685|1913688611^TN^20230328||N||Y|5||||20220327|Y
PV1|1|C
```
This is using the default template found inside main.go.

## Executing mockaroo2hl7v2

There are four possible command line arguments:

```
  -input string
    	Input File.  STDIN used if not specified
  -mllp
    	Wrap each message in MLLP envelope?  Default is false.
  -print0
    	End each message with ASCII NUL?  Default is false.
  -template string
    	Output Template File.  Default can be found in main.go
```

### MLLP

Specifying -mllp when running will wrap the generated message or messages with the _standard_ MLLP envelop.  That is:

 &lt;SB>The Generated HL7 Message...&lt;EB>&lt;CR>

* &lt;SB> = Start Block.  This will be set to ASCII Vertical Tab, hex 0x0B
* &lt;EB> = End Block.  This will be set to ASCII File Separator, hex 0x1C
* &lt;CR> = Carriage Return, hex 0x0D

This _might_ be useful if you are piping generated messages directly to something that will send over TCP/IP but _is not_ doing its own wrapping.

### PRINT0

This is something I added, just knowing that I might pipe generated messages through _xargs_ and to make it easy on myself have generated messages separated by the NUL character.  

This allows you to use the -0 flag with xargs to split things up.  

### Examples

Let's say you have the simple JSON example at the top of this page saved in a file named _demo.json_.  You can convert that to an HL7 message and send to a host and port with this command:

```bash
cat demo.json | go run main.go encounter.go -mllp | netcat -c -vv 10.0.0.50 6661
```

Note we use the -mllp flag to wrap the message with the necessary characters for the downstream MLLP listener.  In this example we are piping the generated message to netcat to facilitate transfer... which isn't going to provide MLLP wrapping on its own.

The output will simple show a number of bytes:

```
Total received bytes: 0
Total sent bytes: 560
```

---

Building on the above, lets say we have an array of JSON objects saved in a file named multiple.json.  We can pipe the generated messages out and have each sent via netcat:

```bash
cat HL7JSON.json | go run main.go encounter.go -mllp -print0 | xargs -0 sh -c 'for i; do echo "$i" | netcat -c -vv 10.0.0.50 6661;done' _
```

Here we use the -mllp flag to wrap the message, but also the -print0 flag to separate each of the generated messages with NUL.

Then we pipe this to xargs, with the -0 flag to specify that the incoming data should be split by the NUL character.  Then each message is piped to netcat and sent to the endpoint.

The output in this case will look something like:

```
[10.0.0.50] 6661 open
Total received bytes: 0
Total sent bytes: 726
[10.0.0.50] 6661 open
Total received bytes: 0
Total sent bytes: 852
[10.0.0.50] 6661 open
Total received bytes: 0
Total sent bytes: 774
```
Our input file had 3 JSON objects, which generated 3 HL7 messages.  So we see 3 different sets of _sent bytes_.  

Now obviously there are other tools to send data over TCP/IP or even MLLP via TCP/IP.  I was just using netcat as an example.

---

Again assuming we have a file with an array of JSON objects, we could use the following command to generate HL7 messages for each JSON object and write them all to a file with:

```bash
cat HL7JSON.json | go run main.go encounter.go > examples.hl7
```

Looking at the contents of the examples.hl7 file you'll see several HL7 messages.  One for each of the objects in the input JSON file.

Now, lets say that we want each generated HL7 message saved to a separate file:

```bash
cat HL7JSON.json | go run main.go encounter.go -print0 | xargs -0 sh -c 'for i; do echo "$i" > example"$RANDOM".hl7;done' _
```

The input file having 3 JSON objects, we will see 3 separate files on disk:

```bash
-rw-r--r--  1 amoody  staff  723 Mar 30 22:09 example21080.hl7
-rw-r--r--  1 amoody  staff  849 Mar 30 22:09 example30613.hl7
-rw-r--r--  1 amoody  staff  771 Mar 30 22:09 example4984.hl7
```

Each file will have one HL7 message in it.

## Templates

To generate HL7 messages the Encounter object (found in [encounter.go](encounter.go)) is passed to a Go [Text Template](https://pkg.go.dev/text/template).

If an external template is not specified (via the -template argument) a default is used.  The default can be found in main.go in the DefaultTemplate function.

Getting familiar with the Encounter object while looking at the default template is the best way to understand how to create your own.

Will add some examples here at some point.

## Mockaroo Setup

mockaroo2hl7v2 expects the input to be JSON of a specific layout. A generated JSON schema is in this repo named [hl7_schema.json](hl7_schema.json).

You could generate the JSON in any way you want.  I use Mockaroo for this, hence the mockaroo2hl7v2 name.

I won't go into specifics about using Mockaroo, but it is a great tool if you need to generate fake data.  If you want to use Mockaroo, sign up for an account and import the file: [HL7JSON.mockaroo.json](HL7JSON.mockaroo.json) located in this repo.  From the Mockaroo Schemas page hit the Restore From Backup button to make that happen.
