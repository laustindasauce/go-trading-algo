# Redis Shortable adding Routines
I'm going to try to speed up execution time by adding go-routines.. not sure if possible with the way the program is set up currently, however..

## Can't add them to the modules since we need to run the sequentially, sadly
What's weird is it does the Nasdaq search so quickly but the shortable takes forever

## Ideas
Something I could do is split the nasdaq assets into fifths then check shortability on the 5 as go-routines

## Update
This was all unnecessary as I was able to fix runtime issues by adjusting the code wihout the go-routines.. Since it takes some time to split the slice into sub-slices it actually ended up taking longer to adjust the slice then run as routine than it does to just run with entire slice.