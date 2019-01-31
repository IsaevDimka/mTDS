__FEATURED TODO__
+ Development testing
+ Production testing

__LAST CHANGES__ (other changes see `docs/changelog.txt`)

> __31.01.2019__
+ Appended Log request option to `settings.ini`
+ Deployment to SG TDS
+ Debug

> __20.01.2019__
+ Testing some cases related c/build/click
+ Documentation for modules and thin places
+ Removin database descriptions when update flows
+ Preparing to release
+ Reformated some parts
+ Single click implementation (viewer)
+ Release 2.0.5alpha

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
