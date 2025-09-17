# SlogColor Enhancement: Source File Colorization and Medium Path Mode

## Overview

This document provides detailed instructions for implementing two enhancements to the slogcolor library:

1. **Source File Colorization**: Add the ability to colorize the source file output using the same pattern as other fields
2. **Medium Path Mode**: Add a new `MediumFile` mode that shows relative paths from the project root

## Implementation Requirements

### 1. Source File Colorization

**Goal**: Allow users to configure a color for the source file information output, following the same pattern used for other colorizable fields in the library.

**Files to Modify**:
- `options.go` - Add new color option field
- `handler.go` - Apply color when formatting source file output

**Detailed Steps**:

#### Step 1.1: Add Color Option to Options Struct

In `options.go`, add a new field to the `Options` struct:

```go
// SrcFileColor is the color of the source file info, default to empty (no color).
SrcFileColor *color.Color
```

#### Step 1.2: Update DefaultOptions

In `options.go`, add the new field to `DefaultOptions`:

```go
var DefaultOptions *Options = &Options{
    // ... existing fields ...
    SrcFileColor:  color.New(), // Default to no color
    // ... rest of fields ...
}
```

#### Step 1.3: Apply Color in Handler

In `handler.go`, modify the source file formatting section (around line 85-95) to apply the color:

**Current code**:
```go
fmt.Fprint(bf, formatted)
```

**New code**:
```go
if h.opts.SrcFileColor == nil {
    h.opts.SrcFileColor = color.New() // set to empty otherwise we have a null pointer
}
fmt.Fprint(bf, h.opts.SrcFileColor.Sprint(formatted))
```

### 2. Medium Path Mode Implementation

**Goal**: Add a new `MediumFile` source file mode that displays the relative path from the project root instead of just the filename or full absolute path.

**Files to Modify**:
- `sourceFileMode.go` - Add new enum constant
- `handler.go` - Implement the new path mode logic

**Detailed Steps**:

#### Step 2.1: Add MediumFile Constant

In `sourceFileMode.go`, add the new constant:

```go
const (
    // Nop does nothing.
    Nop SourceFileMode = iota

    // ShortFile produces only the filename (for example main.go:69).
    ShortFile

    // MediumFile produces the relative file path from project root (for example cmd/server/main.go:69).
    MediumFile

    // LongFile produces the full file path (for example /home/user/go/src/myapp/main.go:69).
    LongFile
)
```

#### Step 2.2: Implement Medium Path Logic

In `handler.go`, modify the filename assignment logic in the `Handle` method (around line 75-85):

**Current code**:
```go
var filename string
switch h.opts.SrcFileMode {
case Nop:
    break
case ShortFile:
    filename = filepath.Base(f.File)
case LongFile:
    filename = f.File
}
```

**New code**:
```go
var filename string
switch h.opts.SrcFileMode {
case Nop:
    break
case ShortFile:
    filename = filepath.Base(f.File)
case MediumFile:
    filename = h.getRelativePath(f.File)
case LongFile:
    filename = f.File
}
```

#### Step 2.3: Add Helper Method for Relative Path

Add this new method to the `Handler` struct in `handler.go`:

```go
// getRelativePath returns the file path relative to the project root
func (h *Handler) getRelativePath(fullPath string) string {
    // Try to get the working directory (project root)
    if wd, err := os.Getwd(); err == nil {
        if relPath, err := filepath.Rel(wd, fullPath); err == nil {
            return relPath
        }
    }
    // Fallback to full path if we can't determine relative path
    return fullPath
}
```

**Note**: You'll need to add `"os"` to the imports in `handler.go` if it's not already there.

#### Step 2.4: Update DefaultOptions (Optional)

Consider updating the default `SrcFileMode` in `options.go` to use the new `MediumFile` mode:

```go
var DefaultOptions *Options = &Options{
    // ... other fields ...
    SrcFileMode:   MediumFile, // Changed from ShortFile
    // ... rest of fields ...
}
```

## Testing the Implementation

### Basic Test Cases

1. **Test Source File Colorization**:
   ```go
   opts := slogcolor.DefaultOptions
   opts.SrcFileColor = color.New(color.FgYellow)
   logger := slog.New(slogcolor.NewHandler(os.Stderr, opts))
   logger.Info("Test message")
   ```

2. **Test Medium Path Mode**:
   ```go
   opts := slogcolor.DefaultOptions
   opts.SrcFileMode = slogcolor.MediumFile
   logger := slog.New(slogcolor.NewHandler(os.Stderr, opts))
   logger.Info("Test message from nested file")
   ```

3. **Test Combined Features**:
   ```go
   opts := slogcolor.DefaultOptions
   opts.SrcFileMode = slogcolor.MediumFile
   opts.SrcFileColor = color.New(color.FgCyan, color.Bold)
   logger := slog.New(slogcolor.NewHandler(os.Stderr, opts))
   logger.Info("Colorized medium path test")
   ```

### Edge Cases to Consider

1. **Relative Path Calculation**: Ensure the relative path calculation works correctly from different working directories
2. **Null Color Pointer**: Ensure the color field is properly initialized to prevent null pointer panics
3. **Path Truncation**: Test that the `SrcFileLength` option still works correctly with the new medium path mode
4. **Go Module Root**: Consider if paths should be relative to the Go module root instead of current working directory

## Documentation Updates

Don't forget to update:

1. **README.md**: Add examples showing the new features
2. **Code comments**: Ensure all new fields and constants have proper documentation comments
3. **Example usage**: Update `example/main.go` to demonstrate the new features

## Implementation Notes

- Follow the existing code style and patterns in the library
- Ensure thread safety is maintained (the existing mutex usage should be sufficient)
- Test thoroughly with different project structures and working directories
- Consider performance impact of the `os.Getwd()` and `filepath.Rel()` calls in the hot path
- Maintain backward compatibility with existing configurations

## Alternative Implementations

For the relative path calculation, you might also consider:

1. **Go Module Root Detection**: Use `go list -m` or similar to find the module root
2. **Configurable Base Path**: Allow users to specify a custom base path for relative calculation
3. **Caching**: Cache the working directory to avoid repeated `os.Getwd()` calls

Choose the approach that best fits the library's design philosophy and performance requirements.