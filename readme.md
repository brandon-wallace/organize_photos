
# Photo Organizer

## A lightweight, fast command-line tool written in Go to automatically organize your photos and RAW files into a clean, date-based directory structure using EXIF metadata.
Features

EXIF-Based Sorting: Automatically reads the DateTime tag from image metadata to sort files.
Standardized Structure: Organizes files into a YYYY/MM/DD hierarchy.
Smart Fallbacks: Images without EXIF data are safely moved to an Unknown folder instead of throwing errors.
Duplicate Prevention: Skips moving files if a file with the same name already exists in the destination, preventing accidentally overwriting files.
Multiple Formats: Supports common JPEG formats (.jpg, .jpeg) and Sony RAW files (.arw) out of the box.

### Prerequisites

Go 1.21 or later (due to the use of log/slog).

## Installation

## Clone and Build

```
git clone https://github.com/brandon-wallace/organize_photos.git

cd organize_photos

go build -o organize_photos
```

## Usage

The tool requires a source directory (-src) where your messy photos live, and a destination directory (-dst) where you want the organized folders to be created.
Basic Example

```
./organize_photos -src /path/to/unsorted/photos -dst /path/to/organized/library
```

## Default Behavior

If you don't provide flags, it defaults to using your current working directory for both source and destination.

```
./organize_photos
```

## Flags

| Flag | Description | Default |
| --- | --- | --- |
| `-src` | Source directory containing images | `.` (Current directory) |
| `-dst` | Destination directory for organized folders | `.` (Current directory) |

## Example Output

Running the command will output a clean summary of what was processed:

```
INFO Organizing image files... source=/home/user/Downloads dst=/home/user/Pictures

Organization complete.
 Successfully moved: 142
 Duplicate files:    5
 Errors:             0
 Total processed:    147
```

## How It Works

Scans: It looks through the -src directory for files matching .jpg, .jpeg, and .arw (case-insensitive).
Reads: It opens each file and attempts to extract the EXIF DateTime original tag.
Formats: It formats that date into YYYY/MM/DD.
Creates: It creates the YYYY/MM/DD folder structure inside the -dst directory (or an Unknown folder if no EXIF data is found).
Moves: It moves the file. If a file with the exact same name already exists in the destination, it is counted as a duplicate and left alone.

## Resulting Directory Structure
 
``` 
dst/
├── 2025/
│   └── 08/
│       └── 15/
│           ├── IMG_001.jpg
│           └── IMG_002.arw
├── 2026/
│   └── 09/
│       └── 01/
│           └── holiday_photo.jpg
└── Unknown/
    └── screenshot.jpg
```
