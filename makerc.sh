#!/bin/sh
git checkout rc
git rebase master
git checkout master
git push origin master --tags
git push origin rc
git describe --long --tags | sed 's/\([^-]*-g\)/r\1/;s/-/./g'