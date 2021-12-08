#!/usr/bin/env bash
mkdir -p E2APtmpFiles
mkdir -p E2SMtmpFiles
cat E2APextFileList.txt | xargs mv -t E2APtmpFiles/.
cat E2SMextFileList.txt | xargs mv -t E2SMtmpFiles/.
cd ../oranE2
rm
find . -type f -name \*.c -exec ../int/autogen/autogen -i --no-top-level-comment -l codev {} \;
find . -type f -name \*.h -exec ../int/autogen/autogen -i --no-top-level-comment -l codev {} \;
mv ../asnTextFiles/E2APtmpFiles/* .
rmdir ../asnTextFiles/E2APtmpFiles
rm converter-example.c
cd ../oranE2SM
find . -type f -name \*.c -exec ../int/autogen/autogen -i --no-top-level-comment -l codev {} \;
find . -type f -name \*.h -exec ../int/autogen/autogen -i --no-top-level-comment -l codev {} \;
mv ../asnTextFiles/E2SMtmpFiles/* .
rmdir ../asnTextFiles/E2SMtmpFiles
rm converter-example.c


