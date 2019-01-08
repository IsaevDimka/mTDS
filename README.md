__FEATURED TODO__
+ Development testing
+ Production testing

__LAST CHANGES__ (other changes see `docs/changelog.txt`)

> __08.01.2019__

+ Improved statistics for average RPS
+ Advanced template based stat output
+ Additional params to watch for stat
+ fx

> __07.01.2018__
+ Improved performance for importing flows
+ Additional statistics /stat, /conf
+ Small refactoring for init.go
+ Additional configurable settings @`settings.ini`
+ Production testing on +10000 rps

__INSTALLATION__

1. clone this repository like this:
> `$ git clone https://github.com/CpaMonstr/metatds.git`

2. in folder `metacds` run following comand
> `$ dep init`

if you have no `dep` installed, just use
> `$ curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh`

3. edit configuration file `settings.ini` to make sure
your settings are correct

4. Run the following command in project folder

for Linux
> `$ ./make.sh` 

or for Windows
> `$ make.cmd`

it will build project and automaticly run it.

Additional information you can find in /docs folder,
there are format specifications and usage examples
