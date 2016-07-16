Quixy Kitty
===========

An pure-Go android game


Compiling
=========

(1) read https://github.com/golang/go/wiki/Mobile

(2) gomobile build -v -o new.apk -target=android/arm github.com/dkolbly/quixykitty

(3) you can build it using X on linux, too, for quick testing!

  sudo apt-get install libegl1-mesa-dev libgles2-mesa-dev libx11-dev
  go install github.com/dkolbly/quixykitty

