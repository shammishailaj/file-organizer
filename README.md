# file-organizer
A small program to organize the contents in a directory of the following format
```
/path/to/directory/filename.extension
```
into directories of the format
```
/path/to/directory/yyyy/mm/dd/filename.extension
```

# Build
---

Navigate to the project root and execute the following command

```
make
```

# Usage
---
To organize files in the current directory

```
./file-organizer
```
To organize files in a specific directory
```
./file-organizer -p /path/to/unorganized/directory
``` 