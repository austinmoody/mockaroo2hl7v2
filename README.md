# mockaroo2hl7v2

This is a simple tool to generate HL7v2 messages from a specific [Mockaroo](https://mockaroo.com) JSON output format.  I needed a way to generate fake HL7v2 messages to test things like Mirth channels + wanted to play around with [Go](https://go.dev).

The output HL7v2 is generating using Go [Text Templates](https://pkg.go.dev/text/template).  You can specify a template file when running, if one is not specified there is one in the main.go file.

## The Input

This utility expects a very specific JSON input, one which I use Mockaroo to setup.  There is a backup of the Mockaroo _schema_ in this repo named TODO.  

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

TODO

## Executing mockaroo2hl7v2

TODO

## Mockaroo Setup

TODO
