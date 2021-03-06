#!/bin/bash

# prepare variable
PRJ="{{ .Project }}"
UNZIP_DIR="{{ .UnZipDir }}"

# unzip
# init new version
cp -r /data/deploy /data/tools/
rm -rf /data/tools/deploy/$PRJ
unzip -d /data/tools/deploy /data/tools/deploy/${PRJ}.zip
[ $UNZIP_DIR != $PRJ ] && mv /data/tools/deploy/$UNZIP_DIR /data/tools/deploy/$PRJ

# calculate diff
cd /data/deploy && for file in `find . -type f | grep -v -w logs`; do md5sum $file; done | sort > /data/tools/base.md5
cd /data/tools/deploy && for file in `find . -type f | grep -v -w logs`; do md5sum $file; done | sort > /data/tools/new.md5
diff /data/tools/base.md5 /data/tools/new.md5 > /data/tools/diff

# replace diff
for del_file in `grep '<' /data/tools/diff | awk -F' ./' '{print $2}'`; do
  echo rm `ls -lh /data/deploy/$del_file`
  rm -rf /data/deploy/$del_file;
done
for add_file in `grep '>' /data/tools/diff | awk -F' ./' '{print $2}'`; do
  mkdir -p `dirname /data/deploy/$add_file`
  mv /data/tools/deploy/$add_file /data/deploy/$add_file
  echo mv `ls -lh /data/deploy/$add_file`
done

# clear
rm -rf /data/tools/deploy
rm -rf /data/deploy/*.zip /data/tools/deploy/*.zip
rm -rf /data/tools/*.md5 /data/tools/diff

# set script
#ln -s /data/deploy/$PRJ-$curdate /data/deploy/$PRJ
chown app:app /data/deploy/$PRJ* -R

rm -rf /data/deploy/*.zip