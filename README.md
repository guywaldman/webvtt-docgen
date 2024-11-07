# webvtt-docgen

Generator for documents based on webvtt (video subtitles) files.

The use-case is when you recorded a meeting (e.g., on Zoom) and have a WebVTT file (`.vtt`) which you want to generate a convenient document for, for making meeting recordings more accessible, especially when people download the videos for offline viewing.

## Usage

```shell
go run main.go -- /path/to/file.vtt /path/to/output.md
```