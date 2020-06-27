# srt-overlap-slicer
Slices overlapping subtitles in SRT that have been converted from ASS subtitles

Stolen from @nimatrueway's Gist (https://gist.github.com/nimatrueway/4589700f49c691e5413c5b2df4d02f4f) and adapted to my needs.

### Problem
I want to watch animes in my Smart TV. Although my TV can read ASS subtitles, it does not render them properly, and overlapping subtitles get ignored.
The same happens for SRTs. So, whenever a two subtitles overlap, a third subtitle, representing the intersection, is created.

### What it do
- split in slices overlapping subtitles
- Removes empty subtitles
- Removes subtitles with very very short duration (<200ms)

### Installation/Usage
1. You need a recent version of `go`
2. Clone repo, `cd` to repo folder, run `go build`
3. This will create an executable named `srt-overlap-slicer`
4. Run that executable on the SRT file.
5. This will back up the original `.srt` with a `.bak` extension, and generate a cleaned SRT file
