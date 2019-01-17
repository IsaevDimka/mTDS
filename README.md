__FEATURED TODO__
+ Development testing
+ Production testing

__LAST CHANGES__ (other changes see `docs/changelog.txt`)

> __16.01.2019__
+ Optimized settings to USE SSL yes|no  
+ Settings to backup clicks yes|no 
+ Preparing to release
+ fx settings

> __15.01.2019__
+ Optimized all flows current/new work with
  new params flow_hash/click_hash
+ Performing to load current clicks to redis for
  direct flows
+ URL reformatting /c /r
+ Preparing to release

> __14.01.2019__
+ Imporved statistics template
+ Optimization statistics
+ Preparing to release

> __13.01.2019__
+ Imporved building a click by parameters
+ Optimization for reading flow infos
+ Reformating some code + fx
+ Preparing to release

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
