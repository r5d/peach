# peach

A weather site.

## building

```bash
make
```

## running

```
peach [ -p PORT ]
```

If the port is not given, it defaults to `8151`.

### environment variables

- `PEACH_PHOTON_URL`: Photon api url. Set this if geocoding should be
  enabled.
