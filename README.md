# htmlformfill

Dynamically set HTML form values with a map of field names and desired values.

## API

````go
  Fill(r io.Reader, f map[string]string) (*bytes.Reader, error)
````

## Supported Inputs

- input
  - text
  - number
  - date
  - radio
  - checkbox
- textarea
- select

## Testing

Tests are run against the `form.html` file in the `examples` directory. The test will generate a `test-output.html` file in the `examples` directory regardless of the test outcome.
