# peach

Barebones [weather service][peach].

[peach]: https://peach.ricketyspace.net

## building

Build requirements:

 - make
 - go version 19 or higher

To build the peach binary, just do:

```bash
make
```

## running

```
peach [ -p PORT ]
```

If the port is not given, it defaults to `8151`.

### environment variables

- `PEACH_PHOTON_URL`: Photon API URL. Set this if geocoding should be
  enabled.
