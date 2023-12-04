# /!bin/bash
rm -rf dev
mkdir dev

set -e

cd dev

git clone https://github.com/JordiSubira/scion.git
cd scion
git checkout dispatcher_off
git checkout 76805a00890201f1e633ed9f9acd38f7ce3243d9

cd ..
git clone https://github.com/martenwallewein/bootstrapper.git
cd bootstrapper
git checkout test/mock-only
cd ..
