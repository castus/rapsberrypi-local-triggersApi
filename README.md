# Trigger Api for Raspberry Pi project

## How it works?

Based on the Sun position API, it calculates what LED color and brightness should be set.

I divided a day into parts with different triggers. Documentation can be found in 
[lightController.go](src%2FlightController%2FlightController.go) file.

## Response

```json
{
  "current": "string of a current trigger id, based on request time",
  "checkpoints": [
    "array of times",
    "when there is a need to refresh LED brightness and color",
    "for ex. Sat Dec 10 17:57:33 CET 2022"
  ]
}
```
