# CoapCloud Demo

This is a simple demo of the basic UDP proxying to a separate coap-capable backend

## Getting Started

`make`, `make all`, or `make run` to boot the udp proxy and the handler funcs

### How does it work?

The calculator server starts at 0.

A POST with `+2` in the payload will add 2

A POST with `-2` in the payload will subtract 2

A GET will return the current calculator total

A DELETE will clear the calculator back to 0 again
