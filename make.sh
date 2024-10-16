#!/bin/sh

rm ./AlfKoreanSearch.alfredworkflow
cd workflow
sed "s/{{VERSION_INFO}}/${GITHUB_REF##*/v}/g" < info.plist > info.plist.bak
mv info.plist.bak info.plist
zip -r ../AlfKoreanSearch.alfredworkflow .
