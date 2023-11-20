# /!bin/bash
rm -rf dev
mkdir dev

set -e

cd dev

git clone https://github.com/JordiSubira/scion.git
cd scion
git checkout dispatcher_off

cd ..
git clone https://github.com/martenwallewein/bootstrapper.git
cd bootstrapper
git checkout test/mock-only
cd ..
