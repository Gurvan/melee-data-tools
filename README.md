# melee-data-tools

Tools and parsers for inspecting Melee data files.

## Fighter Review Tool

`cmd/fighter-review` parses a fighter `.dat` file and opens a small terminal menu for inspecting the parsed data without editing a scratch test.

```sh
go run ./cmd/fighter-review -- /path/to/PlMs.dat
```

The interactive menu lets you inspect sections such as:

- `summary`
- `descriptor`
- `header`
- `roots`
- `relocation`
- `common`
- `special`
- `model-params`
- `hurtboxes`
- `ecb`
- `jostle`
- `items`
- `actions`
- `shield-pose`
- `animations`
- `subactions`
- `model`

For scripted use, list or print sections directly:

```sh
go run ./cmd/fighter-review -list -- /path/to/PlMs.dat
go run ./cmd/fighter-review -section special -- /path/to/PlMs.dat
```

For a fighter file named like `PlMs.dat`, the tool automatically tries to load `PlMsAJ.dat` from the same directory. You can also use `-animations` to provide the fighter animation bundle explicitly:

```sh
go run ./cmd/fighter-review -animations /path/to/PlMsAJ.dat -section animations -- /path/to/PlMs.dat
```

If no animation bundle is found, the `animations` section lists only the animation offsets referenced by actions in the fighter data file.

### Special Attribute Generation

When run in interactive mode, the tool checks whether a lowercase special attribute file exists for the fighter, such as:

```text
fighter/attributes/ms.go
```

If it does not exist, the tool offers to create it.

You can also generate or overwrite the file explicitly:

```sh
go run ./cmd/fighter-review \
  -write-special \
  -special-id Ms \
  -force \
  -- /path/to/PlMs.dat
```

This creates a Go struct with one field per 4-byte special attribute value:

```go
type Ms struct {
	Unk0x00 float32 // 0x0
	Unk0x04 float32 // 0x4
}
```

The generated file also includes a `String()` method using `helpers.PrettyString`, matching the existing special attribute files.

The tool updates `fighter/attributes/special.go` so the fighter root, such as `ftDataMars`, maps to the generated type.

### Type Inference

By default, generation uses:

```sh
-special-type auto
```

Auto mode assumes `float32` unless the raw 4-byte value looks wrong as a float:

- suspicious float plus signed integer in `[-10000, 10000]` -> `int32`
- suspicious float plus clean byte/sentinel pattern like `0xff000000` -> `Hex32`
- otherwise -> `float32`

`Hex32` values render as strings like `"0xff000000"` in pretty output.

You can force one type for all generated fields:

```sh
-special-type float32
-special-type int32
-special-type uint32
-special-type Hex32
```
